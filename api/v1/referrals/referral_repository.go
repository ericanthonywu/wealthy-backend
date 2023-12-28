package referrals

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/referrals/entities"
	"gorm.io/gorm"
)

type (
	ReferralRepository struct {
		db *gorm.DB
	}

	IReferralRepository interface {
		AccountProfile(personalID uuid.UUID) (data entities.ReferralAccountProfile, err error)
		AccountProfileByRefCode(refCode string) (data entities.ReferralAccountProfileRefCode, err error)
		FirstNode(refCode string) (data entities.ReferralUserReward, err error)
		MemberNode(refCode string) (data []entities.ReferralUserReward, err error)
		GetPreviousCommission(refCodeRefference string) (data entities.PreviousCommission, err error)
		SaveWithdraws(model *entities.WithdrawEntities) (id uuid.UUID, err error)
		GetTierReferralCode(referralCode string) (data []entities.GetReferralInfo, err error)
		GetAccountInfoFromRefCode(referralCode string) (data entities.GetAccountInfo, err error)
		GetCommission(refCode string) (data entities.Commission, err error)
	}
)

func NewReferralRepository(db *gorm.DB) *ReferralRepository {
	return &ReferralRepository{db: db}
}

func (r *ReferralRepository) AccountProfile(personalID uuid.UUID) (data entities.ReferralAccountProfile, err error) {
	if err := r.db.Raw(`SELECT tmg.id as id_gender, pa.file_name ,pa.image_path, pa.id,pa.username, pa.name, pa.dob as date_of_birth, pa.refer_code, pa.email, tmat.account_type, tmg.gender_name as gender, tmr.roles as user_roles
FROM tbl_personal_accounts pa
INNER JOIN tbl_master_account_types tmat ON tmat.id = pa.id_master_account_types
LEFT JOIN tbl_master_genders tmg ON tmg.id = pa.id_master_gender
INNER JOIN tbl_authentications ta ON ta.id_personal_accounts = pa.id
INNER JOIN tbl_master_roles tmr ON tmr.id = ta.id_master_roles
WHERE pa.id=?`, personalID).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.ReferralAccountProfile{}, err
	}
	return data, nil
}

func (r *ReferralRepository) AccountProfileByRefCode(refCode string) (data entities.ReferralAccountProfileRefCode, err error) {
	if err := r.db.Raw(`SELECT tpa.username, tpa.name, tmat.account_type FROM tbl_personal_accounts tpa
INNER JOIN tbl_master_account_types tmat ON tpa.id_master_account_types = tmat.id
WHERE tpa.refer_code = ?`, refCode).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.ReferralAccountProfileRefCode{}, err
	}
	return data, nil
}

func (r *ReferralRepository) FirstNode(refCode string) (data entities.ReferralUserReward, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_user_rewards tur WHERE tur.ref_code=? ORDER BY tur.level ASC`, refCode).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.ReferralUserReward{}, err
	}
	return data, nil
}

func (r *ReferralRepository) MemberNode(refCode string) (data []entities.ReferralUserReward, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_user_rewards tur
WHERE tur.ref_code_reference=?
ORDER BY tur.level ASC`, refCode).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.ReferralUserReward{}, err
	}
	return data, nil
}

func (r *ReferralRepository) GetPreviousCommission(refCodeRefference string) (data entities.PreviousCommission, err error) {

	if err := r.db.Raw(`SELECT total_comission FROM tbl_user_rewards WHERE tbl_user_rewards.ref_code=?`, refCodeRefference).Scan(&data).Error; err != nil {
		logrus.Error((err))
		return entities.PreviousCommission{}, err
	}

	return data, nil
}

func (r *ReferralRepository) SaveWithdraws(model *entities.WithdrawEntities) (id uuid.UUID, err error) {
	if err := r.db.Create(&model).Error; err != nil {
		logrus.Error(err.Error())
		return uuid.Nil, err
	}
	return model.ID, nil
}

func (r *ReferralRepository) GetTierReferralCode(referralCode string) (data []entities.GetReferralInfo, err error) {
	if err := r.db.Raw(`WITH RECURSIVE recursiveTable AS (
      SELECT ref_code, ref_code_reference, 0 as level
      FROM tbl_user_rewards
      WHERE ref_code = ?
   UNION ALL
      SELECT tbl_user_rewards.ref_code, tbl_user_rewards.ref_code_reference, recursiveTable.level + 1
      FROM tbl_user_rewards
         JOIN recursiveTable ON tbl_user_rewards.ref_code_reference = recursiveTable.ref_code )
SELECT * FROM recursiveTable`, referralCode).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.GetReferralInfo{}, err
	}
	return data, nil
}

func (r *ReferralRepository) GetAccountInfoFromRefCode(referralCode string) (data entities.GetAccountInfo, err error) {
	if err := r.db.Joins(`JOIN tbl_master_account_types ON tbl_personal_accounts.id_master_account_types = tbl_master_account_types.id`).
		Where("refer_code=?", referralCode).
		Select("tbl_personal_accounts.name", "tbl_master_account_types.account_type").
		Find(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.GetAccountInfo{}, err
	}
	return data, nil
}

func (r *ReferralRepository) GetCommission(refCode string) (data entities.Commission, err error) {
	if err := r.db.Raw(`SELECT total_comission FROM tbl_user_rewards WHERE ref_code=?`, refCode).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.Commission{}, err
	}
	return data, nil
}