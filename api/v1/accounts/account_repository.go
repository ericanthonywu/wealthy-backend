package accounts

import (
	"errors"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type (
	AccountRepository struct {
		db *gorm.DB
	}

	IAccountRepository interface {
		MasterRoles() (masterRoles entities.AccountMasterRoles)
		MasterAccountType() (id uuid.UUID)
		SignUp(personalEntity *entities.AccountSignUpPersonalAccountEntity, authEntity *entities.AccountSignUpAuthenticationsEntity) (role string, httpCode int, errInfo []errorsinfo.Errors)
		SignUpPersonalAccount(model *entities.AccountSignUpPersonalAccountEntity) (IDPersonalAccount uuid.UUID, err error)
		SignUpAuth(model *entities.AccountSignUpAuthenticationsEntity) (err error)
		SignInAuth(model entities.AccountSignInAuthenticationEntity) (data entities.AccountSignInAuthenticationEntity)
		Profile(IDPersonal uuid.UUID) (data entities.AccountProfile)
	}
)

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) MasterRoles() (masterRoles entities.AccountMasterRoles) {
	if err := r.db.Find(&masterRoles, "roles = ?", "USER").Error; err != nil {
		return entities.AccountMasterRoles{}
	}
	return masterRoles
}

func (r *AccountRepository) MasterAccountType() (id uuid.UUID) {
	var masterAccountTypes entities.AccountMasterAccountType
	if err := r.db.Find(&masterAccountTypes, "account_type = ?", "BASIC").Error; err != nil {
		return uuid.Nil
	}
	return masterAccountTypes.ID
}

func (r *AccountRepository) SignUp(personalEntity *entities.AccountSignUpPersonalAccountEntity, authEntity *entities.AccountSignUpAuthenticationsEntity) (role string, httpCode int, errInfo []errorsinfo.Errors) {
	masterRoles := r.MasterRoles()
	IDMasterAccountType := r.MasterAccountType()

	personalEntity.IDMasterAccountType = IDMasterAccountType
	IDPersonalAccount, err := r.SignUpPersonalAccount(personalEntity)
	if err != nil {
		return "", http.StatusUnprocessableEntity, errorsinfo.ErrorWrapper(errInfo, "", err.Error())
	}

	authEntity.IDPersonalAccounts = IDPersonalAccount
	authEntity.IDMasterRoles = masterRoles.ID
	if err = r.SignUpAuth(authEntity); err != nil {
		return "", http.StatusUnprocessableEntity, errorsinfo.ErrorWrapper(errInfo, "", err.Error())
	}

	return masterRoles.Roles, http.StatusOK, errInfo
}

func (r *AccountRepository) SignUpPersonalAccount(model *entities.AccountSignUpPersonalAccountEntity) (IDPersonalAccount uuid.UUID, err error) {
	result := r.db.First(&model, "email", model.Email)
	if model.ID != uuid.Nil {
		return model.ID, errors.New("email already exist on system")
	}

	model.ID, err = uuid.NewUUID()
	if err != nil {
		logrus.Error(err.Error())
	}

	result = r.db.Create(&model)
	if result.RowsAffected == 0 {
		return uuid.Nil, err
	}

	return model.ID, nil
}

func (r *AccountRepository) SignUpAuth(model *entities.AccountSignUpAuthenticationsEntity) (err error) {
	result := r.db.Create(&model)
	if result.RowsAffected == 0 {
		return err
	}
	return nil
}

func (r *AccountRepository) SignInAuth(model entities.AccountSignInAuthenticationEntity) (data entities.AccountSignInAuthenticationEntity) {
	r.db.Raw("SELECT pa.email, a.active, a.password, mr.roles as role FROM tbl_personal_accounts pa "+
		"INNER JOIN tbl_authentications a ON a.id_personal_accounts = pa.id "+
		"INNER JOIN tbl_master_roles mr ON mr.id = a.id_master_roles "+
		"WHERE email= ? AND a.active = true", model.Email).Scan(&data)

	return data
}

func (r *AccountRepository) Profile(IDPersonal uuid.UUID) (data entities.AccountProfile) {
	if err := r.db.Raw(`SELECT pa.id,pa.username, pa.name, pa.dob as date_of_birth, pa.refer_code, tmat.account_type, tmg.gender_name as gender, tmr.roles as user_roles
FROM tbl_personal_accounts pa
INNER JOIN tbl_master_account_types tmat ON tmat.id = pa.id_master_account_types
LEFT JOIN tbl_master_genders tmg ON tmg.id = pa.id_master_gender
INNER JOIN tbl_authentications ta ON ta.id_personal_accounts = pa.id
INNER JOIN tbl_master_roles tmr ON tmr.id = ta.id_master_roles
WHERE pa.id=?`, IDPersonal).Scan(&data).Error; err != nil {
		return entities.AccountProfile{}
	}
	return data
}
