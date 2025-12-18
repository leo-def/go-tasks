package collaborator

import (
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
    database.Repository
    Activate(token string) error
    CheckByIdAndCompanyId(id uuid.UUID, companyId uuid.UUID) (bool, error)
    Create(collaborator *CollaboratorInfo) error
    CreateTx(tx database.Transaction, collaborator *CollaboratorInfo) error
    CreateWithToken(collaborator *CollaboratorInfoWithToken) error
    CreateWithTokenTx(tx database.Transaction, collaborator *CollaboratorInfoWithToken) error
    Delete(id uuid.UUID) (bool, error)
    DeleteTx(tx database.Transaction, id uuid.UUID) (bool, error)
    GetByAccountId(accountId uuid.UUID) ([]AccountCollaborator, int64, error)
    GetByAccountIdAndCompanyId(accountId uuid.UUID, companyId uuid.UUID) (*Collaborator, bool, error)
    GetByCompanyId(p httpx.PaginationParams, companyId uuid.UUID) ([]CompanyCollaborator, int64, error)
    GetByCompanyIdAndRoles(p httpx.PaginationParams, companyId uuid.UUID, roles []CollaboratorRole) ([]CompanyCollaborator, int64, error)
    GetById(id uuid.UUID) (Collaborator, bool, error)
    GetByIdAndAccountId(id uuid.UUID, accountId uuid.UUID) (*Collaborator, bool, error)
    GetCompanyCollaboratorById(id uuid.UUID, companyId uuid.UUID) (*CompanyCollaborator, bool, error)
    GetInfoById(id uuid.UUID) (CollaboratorInfo, bool, error)
    GetRoleById(id uuid.UUID) (CollaboratorRole, bool, error)
    SetActivationTokenTx(tx database.Transaction, id uuid.UUID, token string) error
    Update(collaborator *Collaborator) error
}

type gormRepository struct {
	db database.Connection
}

func NewRepository(conn database.Connection) Repository {
	return &gormRepository{db: conn}
}

func (r *gormRepository) GetConnection() database.Connection {
	return r.db
}

func (r *gormRepository) GetInfoById(id uuid.UUID) (CollaboratorInfo, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var collaborator CollaboratorInfo
	result := db.Where("id = ?", id).
		First(&collaborator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return collaborator, false, nil
		}
		return collaborator, false, result.Error
	}
	return collaborator, true, nil
}

func (r *gormRepository) GetRoleById(id uuid.UUID) (CollaboratorRole, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var roleStr string
	res := db.Raw("SELECT role FROM collaborators WHERE id = ?", id).Scan(&roleStr)
	if res.Error != nil {
		return "", false, res.Error
	}
	if res.RowsAffected == 0 {
		return "", false, nil
	}
	return CollaboratorRole(roleStr), true, nil
}

func (r *gormRepository) CheckByIdAndCompanyId(id uuid.UUID, companyId uuid.UUID) (bool, error) {
	db, _ := database.AsGormConn(r.db)
	var exists bool
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM collaborators WHERE id = ? AND company_id = ?)", id, companyId).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *gormRepository) GetById(id uuid.UUID) (Collaborator, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var collaborator Collaborator
	result := db.Where("id = ?", id).
		Preload("Account").
		Preload("Company").
		First(&collaborator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return collaborator, false, nil
		}
		return collaborator, false, result.Error
	}
	return collaborator, true, nil
}

func (r *gormRepository) GetByAccountId(accountId uuid.UUID) ([]AccountCollaborator, int64, error) {
	db, _ := database.AsGormConn(r.db)
	var collaborators []AccountCollaborator
	var count int64

	q := db.Model(&AccountCollaborator{}).
		Where("account_id = ?", accountId).
		Preload("Company")
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Find(&collaborators).Error; err != nil {
		return nil, 0, err
	}
	return collaborators, count, nil
}

func (r *gormRepository) GetByCompanyId(p httpx.PaginationParams, companyId uuid.UUID) ([]CompanyCollaborator, int64, error) {
	db, _ := database.AsGormConn(r.db)
	return database.Paginate[CompanyCollaborator](db, &CompanyCollaborator{}, p, func(q *gorm.DB) *gorm.DB {
		q = q.Where("collaborators.company_id = ?", companyId)
		q = q.Preload("Account")
		if p.Filter != "" {
			q = q.Joins("JOIN accounts a ON a.id = collaborators.account_id")
			q = q.Where("a.name ILIKE ?", "%"+p.Filter+"%")
		}
		return q
	})
}

func (r *gormRepository) GetByCompanyIdAndRoles(p httpx.PaginationParams, companyId uuid.UUID, roles []CollaboratorRole) ([]CompanyCollaborator, int64, error) {
	db, _ := database.AsGormConn(r.db)
	return database.Paginate[CompanyCollaborator](db, &CompanyCollaborator{}, p, func(q *gorm.DB) *gorm.DB {
		q = q.Where("collaborators.company_id = ?", companyId)
		q = q.Preload("Account")
		if len(roles) == 0 {
			q = q.Where("1 = 0")
		} else {
			q = q.Where("collaborators.role IN ?", roles)
		}
		if p.Filter != "" {
			q = q.Joins("JOIN accounts a ON a.id = collaborators.account_id")
			q = q.Where("a.name ILIKE ?", "%"+p.Filter+"%")
		}
		return q
	})
}

func (r *gormRepository) GetCompanyCollaboratorById(id uuid.UUID, companyId uuid.UUID) (*CompanyCollaborator, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var collaborator CompanyCollaborator
	result := db.Where("id = ? AND company_id = ?", id, companyId).
		Preload("Account").
		First(&collaborator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &collaborator, true, nil
}

func (r *gormRepository) Delete(id uuid.UUID) (bool, error) {
    db, _ := database.AsGormConn(r.db)
    if err := db.Model(&Collaborator{}).Delete(&Collaborator{}, id).Error; err != nil {
        return false, err
    }
    return true, nil
}

func (r *gormRepository) DeleteTx(tx database.Transaction, id uuid.UUID) (bool, error) {
    txdb, _ := database.AsGormTx(tx)
    if err := txdb.Model(&Collaborator{}).Delete(&Collaborator{}, id).Error; err != nil {
        return false, err
    }
    return true, nil
}

func (r *gormRepository) Update(collaborator *Collaborator) error {
	db, _ := database.AsGormConn(r.db)
	q := db.Model(&Collaborator{})
	return q.Where("id = ? AND company_id = ?", collaborator.ID, collaborator.CompanyID).Updates(collaborator).Error
}

func (r *gormRepository) Create(collaborator *CollaboratorInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(collaborator).Error
}

func (r *gormRepository) CreateTx(tx database.Transaction, collaborator *CollaboratorInfo) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(collaborator).Error
}

func (r *gormRepository) CreateWithToken(collaborator *CollaboratorInfoWithToken) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(collaborator).Error
}

func (r *gormRepository) CreateWithTokenTx(tx database.Transaction, collaborator *CollaboratorInfoWithToken) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(collaborator).Error
}

func (r *gormRepository) Activate(token string) error {
	db, _ := database.AsGormConn(r.db)
	q := db.Model(&CollaboratorInfo{}).Where("activation_token = ?", token)

	return q.Updates(map[string]interface{}{
		"active":           true,
		"activation_token": gorm.Expr("NULL"),
	}).Error
}

func (r *gormRepository) SetActivationTokenTx(tx database.Transaction, id uuid.UUID, token string) error {
	txdb, _ := database.AsGormTx(tx)
	q := txdb.Model(&CollaboratorInfo{}).Where("id = ?", id)
	return q.Updates(map[string]interface{}{
		"activation_token": token,
		"active":           true,
	}).Error
}

func (r *gormRepository) GetByAccountIdAndCompanyId(accountId uuid.UUID, companyId uuid.UUID) (*Collaborator, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var collaborator Collaborator
	result := db.Where("account_id = ? AND company_id = ?", accountId, companyId).
		Preload("Account").
		Preload("Company").
		First(&collaborator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &collaborator, true, nil
}

func (r *gormRepository) GetByIdAndAccountId(id uuid.UUID, accountId uuid.UUID) (*Collaborator, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var collaborator Collaborator
	result := db.Where("id = ? AND account_id = ?", id, accountId).
		Preload("Account").
		Preload("Company").
		First(&collaborator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &collaborator, true, nil
}
