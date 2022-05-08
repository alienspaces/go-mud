package prepare

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparable"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
)

// Repository - Methods for preparing and fetching repo statements
type Repository struct {
	Log logger.Logger
	DB  *sqlx.DB
	// prepared
	prepared map[string]bool
	// statements
	getOneStmtList          map[string]*sqlx.Stmt
	getOneForUpdateStmtList map[string]*sqlx.Stmt
	getManyStmtList         map[string]*sqlx.NamedStmt
	createStmtList          map[string]*sqlx.NamedStmt
	updateOneStmtList       map[string]*sqlx.NamedStmt
	updateManyStmtList      map[string]*sqlx.NamedStmt
	deleteOneStmtList       map[string]*sqlx.NamedStmt
	deleteManyStmtList      map[string]*sqlx.NamedStmt
	removeOneStmtList       map[string]*sqlx.NamedStmt
	removeManyStmtList      map[string]*sqlx.NamedStmt

	getOneSQLList          map[string]string
	getOneForUpdateSQLList map[string]string
	getManySQLList         map[string]string
	createSQLList          map[string]string
	updateOneSQLList       map[string]string
	updateManySQLList      map[string]string
	deleteOneSQLList       map[string]string
	deleteManySQLList      map[string]string
	removeOneSQLList       map[string]string
	removeManySQLList      map[string]string
}

var _ preparer.Repository = &Repository{}

// NewPrepare -
func NewRepositoryPreparer(l logger.Logger) (*Repository, error) {

	p := Repository{
		Log: l,

		// prepared
		prepared: make(map[string]bool),

		// statement lists
		getOneStmtList:          make(map[string]*sqlx.Stmt),
		getOneForUpdateStmtList: make(map[string]*sqlx.Stmt),
		getManyStmtList:         make(map[string]*sqlx.NamedStmt),
		createStmtList:          make(map[string]*sqlx.NamedStmt),
		updateOneStmtList:       make(map[string]*sqlx.NamedStmt),
		updateManyStmtList:      make(map[string]*sqlx.NamedStmt),
		deleteOneStmtList:       make(map[string]*sqlx.NamedStmt),
		deleteManyStmtList:      make(map[string]*sqlx.NamedStmt),
		removeOneStmtList:       make(map[string]*sqlx.NamedStmt),
		removeManyStmtList:      make(map[string]*sqlx.NamedStmt),

		getOneSQLList:          make(map[string]string),
		getOneForUpdateSQLList: make(map[string]string),
		getManySQLList:         make(map[string]string),
		createSQLList:          make(map[string]string),
		updateOneSQLList:       make(map[string]string),
		updateManySQLList:      make(map[string]string),
		deleteOneSQLList:       make(map[string]string),
		deleteManySQLList:      make(map[string]string),
		removeOneSQLList:       make(map[string]string),
		removeManySQLList:      make(map[string]string),
	}

	return &p, nil
}

// Init - Initialise preparer with database tx
func (p *Repository) Init(db *sqlx.DB) error {

	if db == nil {
		msg := "database db is nil, cannot init"
		p.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	p.DB = db

	return nil
}

// Prepare - Prepares all repo SQL statements for faster execution
func (p *Repository) Prepare(m preparable.Repository, shouldExclude preparer.ExcludePreparation) error {

	// already prepared
	if _, ok := p.prepared[m.TableName()]; ok {
		return nil
	}

	p.Log.Debug("** Preparing ** table >%s< statements", m.TableName())

	// get one
	if !shouldExclude.GetOne {
		query := m.GetOneSQL()
		getOneStmt, err := p.DB.Preparex(query)
		if err != nil {
			p.Log.Warn("Error preparing GetOneSQL statement >%v<", err)
			return err
		}

		p.getOneSQLList[m.TableName()] = query
		p.getOneStmtList[m.TableName()] = getOneStmt
	}

	// get one for update
	if !shouldExclude.GetOneForUpdate {
		query := m.GetOneForUpdateSQL()

		getOneForUpdateStmt, err := p.DB.Preparex(query)
		if err != nil {
			p.Log.Warn("Error preparing GetOneForUpdateSQL statement >%v<", err)
			return err
		}

		p.getOneForUpdateSQLList[m.TableName()] = query
		p.getOneForUpdateStmtList[m.TableName()] = getOneForUpdateStmt
	}

	// get many
	if !shouldExclude.GetMany {
		query := m.GetManySQL()

		getManyStmt, err := p.DB.PrepareNamed(m.GetManySQL())
		if err != nil {
			p.Log.Warn("Error preparing GetManySQL statement >%v<", err)
			return err
		}

		p.getManySQLList[m.TableName()] = query
		p.getManyStmtList[m.TableName()] = getManyStmt
	}

	// create
	if !shouldExclude.CreateOne {
		query := m.CreateOneSQL()

		createStmt, err := p.DB.PrepareNamed(query)
		if err != nil {
			p.Log.Warn("Error preparing CreateSQL statement >%v<", err)
			return err
		}

		p.createSQLList[m.TableName()] = query
		p.createStmtList[m.TableName()] = createStmt
	}

	// update
	if !shouldExclude.UpdateOne {
		query := m.UpdateOneSQL()

		updateOneStmt, err := p.DB.PrepareNamed(query)
		if err != nil {
			p.Log.Warn("Error preparing UpdateOneSQL statement >%v<", err)
			return err
		}

		p.updateOneSQLList[m.TableName()] = query
		p.updateOneStmtList[m.TableName()] = updateOneStmt
	}

	// update many
	if !shouldExclude.UpdateMany {
		query := m.UpdateManySQL()

		updateManyStmt, err := p.DB.PrepareNamed(query)
		if err != nil {
			p.Log.Warn("Error preparing UpdateManySQL statement >%v<", err)
			return err
		}

		p.updateManySQLList[m.TableName()] = query
		p.updateManyStmtList[m.TableName()] = updateManyStmt
	}

	// delete
	if !shouldExclude.DeleteOne {
		query := m.DeleteOneSQL()

		deleteStmt, err := p.DB.PrepareNamed(query)
		if err != nil {
			p.Log.Warn("Error preparing DeleteSQL statement >%v<", err)
			return err
		}

		p.deleteOneSQLList[m.TableName()] = query
		p.deleteOneStmtList[m.TableName()] = deleteStmt
	}

	// delete many
	if !shouldExclude.DeleteMany {
		query := m.DeleteManySQL()

		deleteManyStmt, err := p.DB.PrepareNamed(query)
		if err != nil {
			p.Log.Warn("Error preparing DeleteManySQL statement >%v<", err)
			return err
		}

		p.deleteManySQLList[m.TableName()] = query
		p.deleteManyStmtList[m.TableName()] = deleteManyStmt
	}

	// remove
	if !shouldExclude.RemoveOne {
		query := m.RemoveOneSQL()

		removeStmt, err := p.DB.PrepareNamed(query)
		if err != nil {
			p.Log.Warn("Error preparing RemoveSQL statement >%v<", err)
			return err
		}

		p.removeOneSQLList[m.TableName()] = query
		p.removeOneStmtList[m.TableName()] = removeStmt
	}

	// remove many
	if !shouldExclude.RemoveMany {
		query := m.RemoveManySQL()

		removeManyStmt, err := p.DB.PrepareNamed(query)
		if err != nil {
			p.Log.Warn("Error preparing RemoveManySQL statement >%v<", err)
			return err
		}

		p.removeManySQLList[m.TableName()] = query
		p.removeManyStmtList[m.TableName()] = removeManyStmt
	}

	p.prepared[m.TableName()] = true

	return nil
}

// GetOneStmt -
func (p *Repository) GetOneStmt(m preparable.Repository) *sqlx.Stmt {

	stmt := p.getOneStmtList[m.TableName()]

	return stmt
}

// GetOneForUpdateStmt -
func (p *Repository) GetOneForUpdateStmt(m preparable.Repository) *sqlx.Stmt {

	stmt := p.getOneForUpdateStmtList[m.TableName()]

	return stmt
}

// GetManyStmt -
func (p *Repository) GetManyStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.getManyStmtList[m.TableName()]

	return stmt
}

// CreateOneStmt -
func (p *Repository) CreateOneStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.createStmtList[m.TableName()]

	return stmt
}

// UpdateOneStmt -
func (p *Repository) UpdateOneStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.updateOneStmtList[m.TableName()]

	return stmt
}

// UpdateManyStmt -
func (p *Repository) UpdateManyStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.updateManyStmtList[m.TableName()]

	return stmt
}

// DeleteOneStmt -
func (p *Repository) DeleteOneStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.deleteOneStmtList[m.TableName()]

	return stmt
}

// DeleteManyStmt -
func (p *Repository) DeleteManyStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.deleteManyStmtList[m.TableName()]

	return stmt
}

// RemoveOneStmt -
func (p *Repository) RemoveOneStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.removeOneStmtList[m.TableName()]

	return stmt
}

// RemoveManyStmt -
func (p *Repository) RemoveManyStmt(m preparable.Repository) *sqlx.NamedStmt {

	stmt := p.removeManyStmtList[m.TableName()]

	return stmt
}

// GetOneSQL -
func (p *Repository) GetOneSQL(m preparable.Repository) string {

	query := p.getOneSQLList[m.TableName()]

	return query
}

// GetManySQL -
func (p *Repository) GetManySQL(m preparable.Repository) string {

	query := p.getManySQLList[m.TableName()]

	return query
}

// CreateSQL -
func (p *Repository) CreateSQL(m preparable.Repository) string {

	query := p.createSQLList[m.TableName()]

	return query
}

// UpdateOneSQL -
func (p *Repository) UpdateOneSQL(m preparable.Repository) string {

	query := p.updateOneSQLList[m.TableName()]

	return query
}

// UpdateManySQL -
func (p *Repository) UpdateManySQL(m preparable.Repository) string {

	query := p.updateManySQLList[m.TableName()]

	return query
}

// DeleteOneSQL -
func (p *Repository) DeleteOneSQL(m preparable.Repository) string {

	query := p.deleteOneSQLList[m.TableName()]

	return query
}

// DeleteManySQL -
func (p *Repository) DeleteManySQL(m preparable.Repository) string {

	query := p.deleteManySQLList[m.TableName()]

	return query
}

// RemoveOneSQL -
func (p *Repository) RemoveOneSQL(m preparable.Repository) string {

	query := p.removeOneSQLList[m.TableName()]

	return query
}

// RemoveManySQL -
func (p *Repository) RemoveManySQL(m preparable.Repository) string {

	query := p.removeManySQLList[m.TableName()]

	return query
}
