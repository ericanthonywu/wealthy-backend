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
		SignUp(personalEntity *entities.AccountSignUpPersonalAccountEntity, authEntity *entities.AccountSignUpAuthenticationsEntity) (role string, idPersonalAccount uuid.UUID, httpCode int, errInfo []errorsinfo.Errors)
		SignUpPersonalAccount(model *entities.AccountSignUpPersonalAccountEntity) (IDPersonalAccount uuid.UUID, err error)
		SignUpAuth(model *entities.AccountSignUpAuthenticationsEntity) (err error)
		SignInAuth(model entities.AccountSignInAuthenticationEntity) (data entities.AccountSignInAuthenticationEntity)
		GetProfile(IDPersonal uuid.UUID) (data entities.AccountProfile)
		GetProfilePassword(IDPersonal uuid.UUID) (data entities.AccountProfilePassword)
		UpdateProfile(customerID uuid.UUID, request map[string]interface{}) (err error)
		UpdatePassword(customerID uuid.UUID, data map[string]interface{}) (err error)
		ForgotPassword(model *entities.AccountForgotPassword) (err error)
		ListRefCode() []string
		GetLevelReferenceCode(referralCode string) (level int, err error)
		WriteRewardsList(model *entities.AccountRewards) (err error)
		WriteOrUpdateReward(model *entities.AccountRewards, idPersonalAccount uuid.UUID) (err error)
		DuplicateExpenseCategory(IDPersonalAccount uuid.UUID) (err error)
		DuplicateExpenseSUbCategory(IDPersonalAccount uuid.UUID) (err error)
		DuplicateIncomeCategory(IDPersonalAccount uuid.UUID) (err error)
		SearchAccount(email string) (data entities.AccountSearchEmail, err error)
		InviteSharing(model *entities.AccountGroupSharing) (err error)
		GetProfileByEmail(email string) (data entities.AccountProfile, err error)
		AcceptSharing(IDSender, IDReceipt uuid.UUID) (err error)
		RejectSharing(IDSender, IDReceipt uuid.UUID) (err error)
		RemoveSharing(IDPASender, IDPARecipient uuid.UUID) (err error)
		RemoveGroupSharingByID(IDGroupSharing uuid.UUID) (err error)
		IDPersonalAccountFromGroupSharing(IDReceiptUUID uuid.UUID) (data entities.AccountPersonalIDGroupSharing)
		GroupSharingInfoByIDPersonalAccount(IDFirstAccount, IDSecondAccount uuid.UUID) (dataFirstAccount entities.AccountGroupSharing, dataSecondAccount entities.AccountGroupSharing)
		GroupSharingList(IDPersonalAccount uuid.UUID) (data []entities.AccountGroupSharingWithProfileInfo, err error)
		ForgotPasswordData(IDPersonalAccount uuid.UUID) (data entities.AccountForgotPassword, err error)
		UpdateForgotPassword(ID uuid.UUID) (err error)
		GenderData(ID uuid.UUID) bool
		IsAlreadySharing(idSender, idRecipient uuid.UUID) bool
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

func (r *AccountRepository) SignUp(personalEntity *entities.AccountSignUpPersonalAccountEntity, authEntity *entities.AccountSignUpAuthenticationsEntity) (role string, idPersonalAccount uuid.UUID, httpCode int, errInfo []errorsinfo.Errors) {
	masterRoles := r.MasterRoles()
	IDMasterAccountType := r.MasterAccountType()

	personalEntity.IDMasterAccountType = IDMasterAccountType
	IDPersonalAccount, err := r.SignUpPersonalAccount(personalEntity)
	if err != nil {
		logrus.Error(err.Error())
		return "", uuid.Nil, http.StatusInternalServerError, errorsinfo.ErrorWrapper(errInfo, "", err.Error())
	}

	authEntity.IDPersonalAccounts = IDPersonalAccount
	authEntity.IDMasterRoles = masterRoles.ID
	if err = r.SignUpAuth(authEntity); err != nil {
		logrus.Error(err.Error())
		return "", uuid.Nil, http.StatusInternalServerError, errorsinfo.ErrorWrapper(errInfo, "", err.Error())
	}

	return masterRoles.Roles, IDPersonalAccount, http.StatusOK, errInfo
}

func (r *AccountRepository) SignUpPersonalAccount(model *entities.AccountSignUpPersonalAccountEntity) (IDPersonalAccount uuid.UUID, err error) {
	_ = r.db.First(&model, "email", model.Email)
	if model.ID != uuid.Nil {
		return model.ID, errors.New("email already exist on system")
	}

	model.ID, err = uuid.NewUUID()
	if err != nil {
		logrus.Error(err.Error())
	}

	if err = r.db.Create(&model).Error; err != nil {
		logrus.Error(err.Error())
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
	r.db.Raw(`SELECT pa.id, pa.email, a.active, a.password, mr.roles as role, tmat.account_type as type
FROM tbl_personal_accounts pa
    INNER JOIN tbl_authentications a ON a.id_personal_accounts = pa.id
    INNER JOIN tbl_master_account_types tmat ON tmat.id = pa.id_master_account_types
INNER JOIN tbl_master_roles mr ON mr.id = a.id_master_roles WHERE email= ? AND a.active = true`, model.Email).Scan(&data)

	return data
}

func (r *AccountRepository) GetProfile(IDPersonal uuid.UUID) (data entities.AccountProfile) {
	if err := r.db.Raw(`SELECT tmg.id as id_gender, pa.file_name ,pa.image_path, pa.id,pa.username, pa.name, pa.dob as date_of_birth, pa.refer_code, pa.email, tmat.account_type, tmg.gender_name as gender, tmr.roles as user_roles, pa.lat, pa.long
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

func (r *AccountRepository) GetProfilePassword(IDPersonal uuid.UUID) (data entities.AccountProfilePassword) {
	if err := r.db.Raw(`SELECT pa.id,pa.username,
       pa.name,
       pa.dob as date_of_birth,
       pa.refer_code,
       pa.email,
       tmat.account_type,
       tmg.gender_name as gender,
       tmr.roles as user_roles,
       ta.password
FROM tbl_personal_accounts pa
INNER JOIN tbl_master_account_types tmat ON tmat.id = pa.id_master_account_types
LEFT JOIN tbl_master_genders tmg ON tmg.id = pa.id_master_gender
INNER JOIN tbl_authentications ta ON ta.id_personal_accounts = pa.id
INNER JOIN tbl_master_roles tmr ON tmr.id = ta.id_master_roles
WHERE pa.id=? AND ta.active=true`, IDPersonal).Scan(&data).Error; err != nil {
		return entities.AccountProfilePassword{}
	}
	return data
}

func (r *AccountRepository) UpdateProfile(customerID uuid.UUID, request map[string]interface{}) (err error) {
	var model entities.AccountSetProfileEntity
	if err := r.db.First(&model, customerID).Error; err != nil {
		return err
	}
	if err := r.db.Model(&model).Updates(request).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) UpdatePassword(customerID uuid.UUID, data map[string]interface{}) (err error) {
	var model entities.AccountAuthorization

	if err := r.db.First(&model, customerID).Error; err != nil {
		return err
	}

	if err := r.db.Model(&model).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) ForgotPassword(model *entities.AccountForgotPassword) (err error) {
	if err = r.db.Create(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *AccountRepository) ListRefCode() []string {
	var referCode []string

	if err := r.db.Model(&entities.AccountSignUpPersonalAccountEntity{}).Distinct("refer_code").Find(&referCode).Error; err != nil {
		logrus.Error(err.Error())
		return []string{}
	}
	return referCode
}

func (r *AccountRepository) GetLevelReferenceCode(referralCode string) (level int, err error) {
	var model entities.AccountRewards

	if err = r.db.Where("ref_code", referralCode).First(&model).Error; err != nil {
		return model.Level, err
	}

	return model.Level, nil
}

func (r *AccountRepository) WriteRewardsList(model *entities.AccountRewards) (err error) {
	if err = r.db.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) WriteOrUpdateReward(modelData *entities.AccountRewards, idPersonalAccount uuid.UUID) (err error) {
	var model entities.AccountRewards

	if err = r.db.First(&model, idPersonalAccount).Error; err != nil {
		return err
	}

	if model.ID != uuid.Nil {
		model.RefCode = modelData.RefCode
		if err = r.db.Save(&model).Error; err != nil {
			return err
		}
	} else {
		if err = r.db.Create(&modelData).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *AccountRepository) DuplicateExpenseCategory(IDPersonalAccount uuid.UUID) (err error) {
	var model interface{}

	if err = r.db.Raw(`INSERT INTO tbl_master_expense_categories_editable (id,expense_types, active, id_personal_accounts,filename, image_path)
SELECT id, expense_types, active, ?, filename, image_path
FROM tbl_master_expense_categories`, IDPersonalAccount).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *AccountRepository) DuplicateExpenseSUbCategory(IDPersonalAccount uuid.UUID) (err error) {
	var model interface{}

	if err = r.db.Raw(`INSERT INTO tbl_master_expense_subcategories_editable (id,subcategories, id_master_expense_categories, active, id_personal_accounts,filename, image_path)
SELECT id, subcategories, id_master_expense_categories,active, ?, filename, image_path
FROM tbl_master_expense_subcategories`, IDPersonalAccount).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *AccountRepository) DuplicateIncomeCategory(IDPersonalAccount uuid.UUID) (err error) {
	var model interface{}

	if err = r.db.Raw(`INSERT INTO tbl_master_income_categories_editable (id,income_types, active, id_personal_accounts,filename, image_path)
SELECT id, income_types,active, ?, filename, image_path
FROM tbl_master_income_categories`, IDPersonalAccount).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *AccountRepository) SearchAccount(email string) (data entities.AccountSearchEmail, err error) {
	if err = r.db.Raw(`SELECT tpa.id FROM tbl_personal_accounts tpa WHERE  tpa.email=?`, email).Scan(&data).Error; err != nil {
		return entities.AccountSearchEmail{}, err
	}
	return data, nil
}

func (r *AccountRepository) InviteSharing(modelGroupSharing *entities.AccountGroupSharing) (err error) {
	if err := r.db.Create(&modelGroupSharing).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) GetProfileByEmail(email string) (data entities.AccountProfile, err error) {
	if err := r.db.Raw(`SELECT tmg.id as id_gender, pa.file_name ,pa.image_path, pa.id,pa.username, pa.name, pa.dob as date_of_birth, pa.refer_code, pa.email, tmat.account_type, tmg.gender_name as gender, tmr.roles as user_roles
FROM tbl_personal_accounts pa
INNER JOIN tbl_master_account_types tmat ON tmat.id = pa.id_master_account_types
LEFT JOIN tbl_master_genders tmg ON tmg.id = pa.id_master_gender
INNER JOIN tbl_authentications ta ON ta.id_personal_accounts = pa.id
INNER JOIN tbl_master_roles tmr ON tmr.id = ta.id_master_roles
WHERE pa.email=?`, email).Scan(&data).Error; err != nil {
		return entities.AccountProfile{}, err
	}
	return data, nil
}

func (r *AccountRepository) AcceptSharing(IDSender, IDReceipt uuid.UUID) (err error) {
	var model interface{}

	// update sender
	if err = r.db.Raw(`UPDATE tbl_group_sharing SET is_accepted=true WHERE id=?`, IDSender).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}

	// update receiver
	if err = r.db.Raw(`UPDATE tbl_group_sharing SET is_accepted=true WHERE id=?`, IDReceipt).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}

	return nil
}

func (r *AccountRepository) RejectSharing(IDSender, IDReceipt uuid.UUID) (err error) {
	var model interface{}

	// delete sender
	if err = r.db.Raw(`DELETE FROM tbl_group_sharing WHERE id=?`, IDSender).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}

	// delete receiver
	if err = r.db.Raw(`DELETE FROM tbl_group_sharing WHERE id=?`, IDReceipt).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *AccountRepository) RemoveSharing(IDPASender, IDPARecipient uuid.UUID) (err error) {
	var model interface {
	}
	// delete sender
	if err = r.db.Raw(`DELETE FROM tbl_group_sharing WHERE id_personal_accounts_share_from=? AND id_personal_accounts_share_to=?`, IDPASender, IDPARecipient).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}

	// delete receiver
	if err = r.db.Raw(`DELETE FROM tbl_group_sharing WHERE id_personal_accounts_share_from=? AND id_personal_accounts_share_to=?`, IDPARecipient, IDPASender).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *AccountRepository) RemoveGroupSharingByID(IDGroupSharing uuid.UUID) (err error) {
	var model interface{}

	// remove group sharing
	if err = r.db.Raw(`DELETE FROM tbl_group_sharing WHERE id=?`, IDGroupSharing).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}

	return nil
}

func (r *AccountRepository) IDPersonalAccountFromGroupSharing(IDReceiptUUID uuid.UUID) (data entities.AccountPersonalIDGroupSharing) {
	if err := r.db.Raw(`SELECT tpa.id as id_personal_accounts, tgs.is_accepted, tpa.name
FROM tbl_group_sharing tgs INNER JOIN tbl_personal_accounts tpa ON tpa.id= tgs.id_personal_accounts_share_from
WHERE tgs.id=?`, IDReceiptUUID).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.AccountPersonalIDGroupSharing{}
	}
	return data
}

func (r *AccountRepository) GroupSharingInfoByIDPersonalAccount(IDFirstAccount, IDSecondAccount uuid.UUID) (dataFirstAccount entities.AccountGroupSharing, dataSecondAccount entities.AccountGroupSharing) {
	if err := r.db.Raw(`SELECT * FROM tbl_group_sharing tgs WHERE tgs.id_personal_accounts_share_from = ? AND tgs.id_personal_accounts_share_to = ?`, IDFirstAccount, IDSecondAccount).Scan(&dataFirstAccount).Error; err != nil {
		logrus.Error(err.Error())
		return entities.AccountGroupSharing{}, entities.AccountGroupSharing{}
	}

	if err := r.db.Raw(`SELECT * FROM tbl_group_sharing tgs WHERE tgs.id_personal_accounts_share_to = ? AND tgs.id_personal_accounts_share_from = ?`, IDFirstAccount, IDSecondAccount).Scan(&dataSecondAccount).Error; err != nil {
		logrus.Error(err.Error())
		return entities.AccountGroupSharing{}, entities.AccountGroupSharing{}
	}

	return dataFirstAccount, dataSecondAccount
}

func (r *AccountRepository) GroupSharingList(IDPersonalAccount uuid.UUID) (data []entities.AccountGroupSharingWithProfileInfo, err error) {
	if err := r.db.Raw(`SELECT tpa.email,
       tmat.account_type as type,
       tpa.file_name     as file_name,
       tpa.image_path    as image_path,
       CASE
           WHEN tgs.is_accepted = false THEN 'pending'
           ELSE 'accepted'
           END           AS status
FROM tbl_group_sharing tgs
         INNER JOIN tbl_personal_accounts tpa ON tgs.id_personal_accounts_share_to = tpa.id
         INNER JOIN tbl_master_account_types tmat ON tmat.id = tpa.id_master_account_types
WHERE tgs.id_personal_accounts_share_from = ?`, IDPersonalAccount).Scan(&data).Error; err != nil {
		return []entities.AccountGroupSharingWithProfileInfo{}, err
	}
	return data, nil
}

func (r *AccountRepository) ForgotPasswordData(IDPersonalAccount uuid.UUID) (data entities.AccountForgotPassword, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_forgot_password tfp
         WHERE tfp.id_personal_accounts=? ORDER BY tfp.created_at DESC LIMIT 1`, IDPersonalAccount).Scan(&data).Error; err != nil {
		return entities.AccountForgotPassword{}, err
	}

	return data, nil
}

func (r *AccountRepository) UpdateForgotPassword(ID uuid.UUID) (err error) {
	var model interface{}

	if err = r.db.Raw(`UPDATE tbl_forgot_password SET is_verified=true WHERE id=?`, ID).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *AccountRepository) GenderData(ID uuid.UUID) bool {
	var data entities.AccountGender

	if err := r.db.Raw(`SELECT EXISTS ( SELECT 1 FROM tbl_master_genders tmg WHERE tmg.id=?)`, ID).Scan(&data).Error; err != nil {
		return data.Exists
	}
	return data.Exists
}

func (r *AccountRepository) IsAlreadySharing(idSender, idRecipient uuid.UUID) bool {
	var data entities.AccountAlreadySharing

	if err := r.db.Raw(`SELECT EXISTS (SELECT 1 FROM tbl_group_sharing tgs WHERE tgs.id_personal_accounts_share_from=? AND tgs.id_personal_accounts_share_to=?)`, idSender, idRecipient).
		Scan(&data).Error; err != nil {
		return data.Exists
	}
	return data.Exists
}