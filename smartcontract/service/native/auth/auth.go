

package auth

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/mixbee/mixbee/account"
	"github.com/mixbee/mixbee/common"
	"github.com/mixbee/mixbee/common/log"
	"github.com/mixbee/mixbee/common/serialization"
	"github.com/mixbee/mixbee/errors"
	"github.com/mixbee/mixbee/smartcontract/service/native"
	"github.com/mixbee/mixbee/smartcontract/service/native/utils"
)

var (
	future = time.Date(2100, 1, 1, 12, 0, 0, 0, time.UTC)
)

func Init() {
	native.Contracts[utils.AuthContractAddress] = RegisterAuthContract
}

/*
 * contract admin management
 */
func initContractAdmin(native *native.NativeService, contractAddr common.Address, mixID []byte) (bool, error) {
	admin, err := getContractAdmin(native, contractAddr)
	if err != nil {
		return false, err
	}
	if admin != nil {
		//admin is already set, just return
		log.Debugf("admin of contract %s is already set", contractAddr.ToHexString())
		return false, nil
	}
	err = putContractAdmin(native, contractAddr, mixID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func InitContractAdmin(native *native.NativeService) ([]byte, error) {
	param := new(InitContractAdminParam)
	rd := bytes.NewReader(native.Input)
	if err := param.Deserialize(rd); err != nil {
		return nil, fmt.Errorf("[initContractAdmin] deserialize param failed: %v", err)
	}
	cxt := native.ContextRef.CallingContext()
	if cxt == nil {
		return nil, fmt.Errorf("[initContractAdmin] no calling context")
	}
	invokeAddr := cxt.ContractAddress

	if !account.VerifyID(string(param.AdminMixID)) {
		return nil, fmt.Errorf("[initContractAdmin] invalid param: adminMixID is %x", param.AdminMixID)
	}
	ret, err := initContractAdmin(native, invokeAddr, param.AdminMixID)
	if err != nil {
		return nil, fmt.Errorf("[initContractAdmin] init failed: %v", err)
	}
	if !ret {
		return utils.BYTE_FALSE, nil
	}

	msg := []interface{}{"initContractAdmin", invokeAddr.ToHexString(), string(param.AdminMixID)}
	pushEvent(native, msg)
	return utils.BYTE_TRUE, nil
}

func transfer(native *native.NativeService, contractAddr common.Address, newAdminMixID []byte, keyNo uint64) (bool, error) {
	admin, err := getContractAdmin(native, contractAddr)
	if err != nil {
		return false, fmt.Errorf("getContractAdmin failed: %v", err)
	}
	if admin == nil {
		log.Debugf("admin of contract %s is not set", contractAddr.ToHexString())
		return false, nil
	}

	ret, err := verifySig(native, admin, keyNo)
	if err != nil {
		return false, fmt.Errorf("verifySig failed: %v", err)
	}
	if !ret {
		log.Debugf("verify Admin's signature failed: admin=%s, keyNo=%d", string(admin), keyNo)
		return false, nil
	}

	adminKey := concatContractAdminKey(native, contractAddr)
	utils.PutBytes(native, adminKey, newAdminMixID)
	return true, nil
}

func Transfer(native *native.NativeService) ([]byte, error) {
	//deserialize param
	param := new(TransferParam)
	rd := bytes.NewReader(native.Input)
	err := param.Deserialize(rd)
	if err != nil {
		return nil, fmt.Errorf("[transfer] deserialize param failed: %v", err)
	}

	if !account.VerifyID(string(param.NewAdminMixID)) {
		return nil, fmt.Errorf("[transfer] invalid param: newAdminMixID is %x", param.NewAdminMixID)
	}
	//prepare event msg
	contract := param.ContractAddr.ToHexString()
	failState := []interface{}{"transfer", contract, false}
	sucState := []interface{}{"transfer", contract, true}

	//call transfer func
	ret, err := transfer(native, param.ContractAddr, param.NewAdminMixID, param.KeyNo)
	if err != nil {
		return nil, fmt.Errorf("[transfer] transfer failed: %v", err)
	}
	if ret {
		pushEvent(native, sucState)
		return utils.BYTE_TRUE, nil
	} else {
		pushEvent(native, failState)
		return utils.BYTE_FALSE, nil
	}
}

func AssignFuncsToRole(native *native.NativeService) ([]byte, error) {
	//deserialize input param
	param := new(FuncsToRoleParam)
	rd := bytes.NewReader(native.Input)
	if err := param.Deserialize(rd); err != nil {
		return nil, fmt.Errorf("[assignFuncsToRole] deserialize param failed: %v", err)
	}

	//prepare event msg
	contract := param.ContractAddr.ToHexString()
	failState := []interface{}{"assignFuncsToRole", contract, false}
	sucState := []interface{}{"assignFuncsToRole", contract, true}

	if param.Role == nil {
		return nil, fmt.Errorf("[assignFuncsToRole] invalid param: role is nil")
	}

	//check the caller's permission
	admin, err := getContractAdmin(native, param.ContractAddr)
	if err != nil {
		return nil, fmt.Errorf("[assignFuncsToRole] getContractAdmin failed: %v", err)
	}
	if admin == nil { //admin has not been set
		return nil, fmt.Errorf("[assignFuncsToRole] admin of contract %s has not been set",
			param.ContractAddr.ToHexString())
	}
	if bytes.Compare(admin, param.AdminMixID) != 0 {
		log.Debugf("[assignFuncsToRole] invalid param: adminMixID doesn't match %s != %s",
			string(param.AdminMixID), string(admin))
		pushEvent(native, failState)
		return utils.BYTE_FALSE, nil
	}
	ret, err := verifySig(native, param.AdminMixID, param.KeyNo)
	if err != nil {
		return nil, fmt.Errorf("[assignFuncsToRole] verify admin's signature failed: %v", err)
	}
	if !ret {
		log.Debugf("[assignFuncsToRole] verifySig return false: adminMixID=%s, keyNo=%d",
			string(admin), param.KeyNo)
		pushEvent(native, failState)
		return utils.BYTE_FALSE, nil
	}

	funcs, err := getRoleFunc(native, param.ContractAddr, param.Role)
	if err != nil {
		return nil, fmt.Errorf("[assignFuncsToRole] getRoleFunc failed: %v", err)
	}
	if funcs != nil {
		funcNames := append(funcs.funcNames, param.FuncNames...)
		funcs.funcNames = stringSliceUniq(funcNames)
	} else {
		funcs = new(roleFuncs)
		funcs.funcNames = stringSliceUniq(param.FuncNames)
	}
	err = putRoleFunc(native, param.ContractAddr, param.Role, funcs)
	if err != nil {
		return nil, fmt.Errorf("[assignFuncsToRole] putRoleFunc failed: %v", err)
	}

	pushEvent(native, sucState)
	return utils.BYTE_TRUE, nil
}

func assignToRole(native *native.NativeService, param *MixIDsToRoleParam) (bool, error) {
	//check admin's permission
	admin, err := getContractAdmin(native, param.ContractAddr)
	if err != nil {
		return false, fmt.Errorf("getContractAdmin failed: %v", err)
	}
	if admin == nil {
		return false, fmt.Errorf("admin of contract %s is not set", param.ContractAddr.ToHexString())
	}
	if bytes.Compare(admin, param.AdminMixID) != 0 {
		log.Debugf("param's adminMixID doesn't match: %s != %s", string(param.AdminMixID),
			string(admin))
		return false, nil
	}
	valid, err := verifySig(native, param.AdminMixID, param.KeyNo)
	if err != nil {
		return false, fmt.Errorf("verify admin's signature failed: %v", err)
	}
	if !valid {
		log.Debugf("[assignMixIDsToRole] verifySig return false: adminMixID=%s, keyNo=%d",
			string(admin), param.KeyNo)
		return false, nil
	}

	//init a permanent auth token
	token := new(AuthToken)
	token.expireTime = uint32(future.Unix())
	token.level = 2
	token.role = param.Role

	for _, p := range param.Persons {
		if p == nil {
			continue
		}
		tokens, err := getMixIDToken(native, param.ContractAddr, p)
		if err != nil {
			return false, fmt.Errorf("getMixIDToken failed: %v", err)
		}
		if tokens == nil {
			tokens = new(roleTokens)
			tokens.tokens = make([]*AuthToken, 1)
			tokens.tokens[0] = token
		} else {
			ret, err := hasRole(native, param.ContractAddr, p, param.Role)
			if err != nil {
				return false, fmt.Errorf("check if %s has role %s failed: %v", string(p),
					string(param.Role), err)
			}
			if !ret {
				tokens.tokens = append(tokens.tokens, token)
			} else {
				continue
			}
		}
		err = putMixIDToken(native, param.ContractAddr, p, tokens)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func AssignMixIDsToRole(native *native.NativeService) ([]byte, error) {
	//deserialize param
	param := new(MixIDsToRoleParam)
	rd := bytes.NewReader(native.Input)
	if err := param.Deserialize(rd); err != nil {
		return nil, fmt.Errorf("[assignMixIDsToRole] deserialize param failed: %v", err)
	}

	if param.Role == nil {
		return nil, fmt.Errorf("[assignMixIDsToRole] invalid param: role is nil")
	}
	for i, mixID := range param.Persons {
		if !account.VerifyID(string(mixID)) {
			return nil, fmt.Errorf("[assignMixIDsToRole] invalid param: param.Persons[%d]=%s",
				i, string(mixID))
		}
	}

	ret, err := assignToRole(native, param)
	if err != nil {
		return nil, fmt.Errorf("[assignMixIDsToRole] failed: %v", err)
	}

	contract := param.ContractAddr.ToHexString()
	failState := []interface{}{"assignMixIDsToRole", contract, false}
	sucState := []interface{}{"assignMixIDsToRole", contract, true}
	if ret {
		pushEvent(native, sucState)
		return utils.BYTE_TRUE, nil
	} else {
		pushEvent(native, failState)
		return utils.BYTE_FALSE, nil
	}
}

func getAuthToken(native *native.NativeService, contractAddr common.Address, mixID, role []byte) (*AuthToken, error) {
	tokens, err := getMixIDToken(native, contractAddr, mixID)
	if err != nil {
		return nil, fmt.Errorf("get token failed, caused by %v", err)
	}
	if tokens != nil {
		for _, token := range tokens.tokens {
			if bytes.Compare(token.role, role) == 0 { //permanent token
				return token, nil
			}
		}
	}
	status, err := getDelegateStatus(native, contractAddr, mixID)
	if err != nil {
		return nil, fmt.Errorf("get delegate status failed, caused by %v", err)
	}
	if status != nil {
		for _, s := range status.status {
			if bytes.Compare(s.role, role) == 0 && native.Time < s.expireTime { //temporary token
				token := new(AuthToken)
				token.role = s.role
				token.level = s.level
				token.expireTime = s.expireTime
				return token, nil
			}
		}
	}
	return nil, nil
}

func hasRole(native *native.NativeService, contractAddr common.Address, mixID, role []byte) (bool, error) {
	token, err := getAuthToken(native, contractAddr, mixID, role)
	if err != nil {
		return false, err
	}
	if token == nil {
		return false, nil
	}
	return true, nil
}

func getLevel(native *native.NativeService, contractAddr common.Address, mixID, role []byte) (uint8, error) {
	token, err := getAuthToken(native, contractAddr, mixID, role)
	if err != nil {
		return 0, err
	}
	if token == nil {
		return 0, nil
	}
	return token.level, nil
}

/*
 * if 'from' has the authority and 'to' has not been authorized 'role',
 * then make changes to storage as follows:
 */
func delegate(native *native.NativeService, contractAddr common.Address, from []byte, to []byte,
	role []byte, period uint32, level uint8, keyNo uint64) (bool, error) {
	var fromHasRole, toHasRole bool
	var fromLevel uint8
	var fromExpireTime uint32

	//check input param
	expireTime := native.Time
	if period+expireTime < period {
		//invalid period param, causing overflow
		return false, fmt.Errorf("[delegate] invalid param: overflow, period=%d", period)
	}
	expireTime = expireTime + period

	//check from's permission
	ret, err := verifySig(native, from, keyNo)
	if err != nil {
		return false, fmt.Errorf("verify %s's signature failed: %v", string(from), err)
	}
	if !ret {
		log.Debugf("verifySig return false: from=%s, keyNo=%d", string(from), keyNo)
		return false, nil
	}

	if !account.VerifyID(string(to)) {
		return false, fmt.Errorf("can not pass MixID validity test: to=%s", string(to))
	}

	//get from's auth token
	fromToken, err := getAuthToken(native, contractAddr, from, role)
	if err != nil {
		return false, fmt.Errorf("getAuthToken of %s failed: %v", string(from), err)
	}
	if fromToken == nil {
		fromHasRole = false
		fromLevel = 0
	} else {
		fromHasRole = true
		fromLevel = fromToken.level
		fromExpireTime = fromToken.expireTime
	}

	//get to's auth token
	toToken, err := getAuthToken(native, contractAddr, to, role)
	if err != nil {
		return false, fmt.Errorf("getAuthToken of %s failed: %v", string(to), err)
	}
	if toToken == nil {
		toHasRole = false
	} else {
		toHasRole = true
	}
	if !fromHasRole || toHasRole {
		log.Debugf("%s doesn't have role %s or %s already has role %s", string(from), string(role),
			string(to), string(role))
		return false, nil
	}

	//check if 'from' has the permission to delegate
	if fromLevel == 2 {
		if level < fromLevel && level > 0 && expireTime < fromExpireTime {
			status, err := getDelegateStatus(native, contractAddr, to)
			if err != nil {
				return false, fmt.Errorf("getDelegateStatus failed: %v", err)
			}
			if status == nil {
				status = new(Status)
			}
			j := -1
			for i, s := range status.status {
				if bytes.Compare(s.role, role) == 0 {
					j = i
					break
				}
			}
			if j < 0 {
				newStatus := &DelegateStatus{
					root: from,
				}
				newStatus.expireTime = expireTime
				newStatus.role = role
				newStatus.level = uint8(level)
				status.status = append(status.status, newStatus)
			} else {
				status.status[j].level = uint8(level)
				status.status[j].expireTime = expireTime
				status.status[j].root = from
			}
			err = putDelegateStatus(native, contractAddr, to, status)
			if err != nil {
				return false, fmt.Errorf("putDelegateStatus failed: %v", err)
			}
			return true, nil
		}
	}
	//TODO: for fromLevel > 2 case
	return false, nil
}

func Delegate(native *native.NativeService) ([]byte, error) {
	//deserialize param
	param := &DelegateParam{}
	rd := bytes.NewReader(native.Input)
	err := param.Deserialize(rd)
	if err != nil {
		return nil, fmt.Errorf("[delegate] deserialize param failed: %v", err)
	}
	if param.Period > 1<<32 || param.Level > 1<<8 {
		return nil, fmt.Errorf("[delegate] period or level is too large")
	}

	//prepare event msg
	contract := param.ContractAddr.ToHexString()
	failState := []interface{}{"delegate", contract, param.From, param.To, false}
	sucState := []interface{}{"delegate", contract, param.From, param.To, true}

	//call the delegate func
	ret, err := delegate(native, param.ContractAddr, param.From, param.To, param.Role,
		uint32(param.Period), uint8(param.Level), param.KeyNo)
	if err != nil {
		return nil, fmt.Errorf("[delegate] failed: %v", err)
	}
	if ret {
		pushEvent(native, sucState)
		return utils.BYTE_TRUE, nil
	} else {
		pushEvent(native, failState)
		return utils.BYTE_FALSE, nil
	}
}

func withdraw(native *native.NativeService, contractAddr common.Address, initiator []byte, delegate []byte,
	role []byte, keyNo uint64) (bool, error) {
	//check from's permission
	ret, err := verifySig(native, initiator, keyNo)
	if err != nil {
		return false, fmt.Errorf("verifySig failed: %v", err)
	}
	if !ret {
		log.Debugf("verifySig return false: initiator=%s, keyNo=%d", string(initiator), keyNo)
		return false, nil
	}

	//code below only works in the case that initiator's level is 2
	//TODO: remove the above limitation
	initToken, err := getAuthToken(native, contractAddr, initiator, role)
	if err != nil {
		return false, fmt.Errorf("getAuthToken failed: %v", err)
	}
	if initToken == nil {
		//initiator does not have the right to withdraw
		log.Debugf("[withdraw] initiator %s does not have the right to withdraw", string(initiator))
		return false, nil
	}
	status, err := getDelegateStatus(native, contractAddr, delegate)
	if err != nil {
		return false, fmt.Errorf("getDelegateStatus failed: %v", err)
	}
	if status == nil {
		return false, nil
	}
	for i, s := range status.status {
		if bytes.Compare(s.role, role) == 0 &&
			bytes.Compare(s.root, initiator) == 0 {
			newStatus := new(Status)
			newStatus.status = append(status.status[:i], status.status[i+1:]...)
			err = putDelegateStatus(native, contractAddr, delegate, newStatus)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

func Withdraw(native *native.NativeService) ([]byte, error) {
	//deserialize param
	param := &WithdrawParam{}
	rd := bytes.NewReader(native.Input)
	err := param.Deserialize(rd)
	if err != nil {
		return nil, fmt.Errorf("[withdraw] deserialize param failed: %v", err)
	}

	//prepare event msg
	contract := param.ContractAddr.ToHexString()
	failState := []interface{}{"withdraw", contract, param.Initiator, param.Delegate, false}
	sucState := []interface{}{"withdraw", contract, param.Initiator, param.Delegate, true}

	//call the withdraw func
	ret, err := withdraw(native, param.ContractAddr, param.Initiator, param.Delegate, param.Role, param.KeyNo)
	if err != nil {
		return nil, fmt.Errorf("[withdraw] withdraw failed: %v", err)
	}
	if ret {
		pushEvent(native, sucState)
		return utils.BYTE_TRUE, nil
	} else {
		pushEvent(native, failState)
		return utils.BYTE_FALSE, nil
	}
}

func verifyToken(native *native.NativeService, contractAddr common.Address, caller []byte, fn string, keyNo uint64) (bool, error) {
	//check caller's identity
	ret, err := verifySig(native, caller, keyNo)
	if err != nil {
		return false, fmt.Errorf("verifySig failed: %v", err)
	}
	if !ret {
		log.Debugf("verifySig return false: caller=%s, keyNo=%d", string(caller), keyNo)
		return false, nil
	}

	//check if caller has the permanent auth token
	tokens, err := getMixIDToken(native, contractAddr, caller)
	if err != nil {
		return false, fmt.Errorf("getMixIDToken failed: %v", err)
	}
	if tokens != nil {
		for _, token := range tokens.tokens {
			funcs, err := getRoleFunc(native, contractAddr, token.role)
			if err != nil {
				return false, fmt.Errorf("getRoleFunc failed: %v", err)
			}
			if funcs == nil || token.expireTime < native.Time {
				continue
			}
			for _, f := range funcs.funcNames {
				if strings.Compare(fn, f) == 0 {
					return true, nil
				}
			}
		}
	}

	status, err := getDelegateStatus(native, contractAddr, caller)
	if err != nil {
		return false, fmt.Errorf("getDelegateStatus failed: %v", err)
	}
	if status != nil {
		for _, s := range status.status {
			funcs, err := getRoleFunc(native, contractAddr, s.role)
			if err != nil {
				return false, fmt.Errorf("getRoleFunc failed: %v", err)
			}
			if funcs == nil || s.expireTime < native.Time {
				continue
			}
			for _, f := range funcs.funcNames {
				if strings.Compare(fn, f) == 0 {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func VerifyToken(native *native.NativeService) ([]byte, error) {
	//deserialize param
	param := &VerifyTokenParam{}
	rd := bytes.NewReader(native.Input)
	err := param.Deserialize(rd)
	if err != nil {
		return nil, fmt.Errorf("[verifyToken] deserialize param failed: %v", err)
	}

	contract := param.ContractAddr.ToHexString()
	failState := []interface{}{"verifyToken", contract, param.Caller, param.Fn, false}
	sucState := []interface{}{"verifyToken", contract, param.Caller, param.Fn, true}

	ret, err := verifyToken(native, param.ContractAddr, param.Caller, param.Fn, param.KeyNo)
	if err != nil {
		return nil, fmt.Errorf("[verifyToken] verifyToken failed: %v", err)
	}
	if !ret {
		pushEvent(native, failState)
		return utils.BYTE_FALSE, nil
	}
	pushEvent(native, sucState)
	return utils.BYTE_TRUE, nil
}

func verifySig(native *native.NativeService, mixID []byte, keyNo uint64) (bool, error) {
	bf := new(bytes.Buffer)
	if err := serialization.WriteVarBytes(bf, mixID); err != nil {
		return false, err
	}
	if err := utils.WriteVarUint(bf, keyNo); err != nil {
		return false, err
	}
	args := bf.Bytes()
	ret, err := native.NativeCall(utils.MixIDContractAddress, "verifySignature", args)
	if err != nil {
		return false, err
	}
	valid, ok := ret.([]byte)
	if !ok {
		return false, errors.NewErr("verifySignature return non-bool value")
	}
	if bytes.Compare(valid, utils.BYTE_TRUE) == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func RegisterAuthContract(native *native.NativeService) {
	native.Register("initContractAdmin", InitContractAdmin)
	native.Register("assignFuncsToRole", AssignFuncsToRole)
	native.Register("delegate", Delegate)
	native.Register("withdraw", Withdraw)
	native.Register("assignMixIDsToRole", AssignMixIDsToRole)
	native.Register("verifyToken", VerifyToken)
	native.Register("transfer", Transfer)
}
