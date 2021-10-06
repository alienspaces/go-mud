package prepare

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparable"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
)

// Prepare - Methods for preparing and fetching repo statements
type Prepare struct {
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
}

var _ preparer.Preparer = &Prepare{}

var getOneSQLList = make(map[string]string)
var getOneForUpdateSQLList = make(map[string]string)
var getManySQLList = make(map[string]string)
var createSQLList = make(map[string]string)
var updateOneSQLList = make(map[string]string)
var updateManySQLList = make(map[string]string)
var deleteOneSQLList = make(map[string]string)
var deleteManySQLList = make(map[string]string)
var removeOneSQLList = make(map[string]string)
var removeManySQLList = make(map[string]string)

// mutex
var mutex = &sync.Mutex{}

// NewPrepare -
func NewPrepare(l logger.Logger) (*Prepare, error) {

	p := Prepare{
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
	}

	return &p, nil
}

// Init - Initialise preparer with database tx
func (p *Prepare) Init(db *sqlx.DB) error {

	if db == nil {
		msg := "Database db is nil, cannot init"
		p.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	p.DB = db

	return nil
}

// Prepare - Prepares all repo SQL statements for faster execution
func (p *Prepare) Prepare(m preparable.Preparable) error {

	// lock/unlock
	mutex.Lock()
	defer mutex.Unlock()

	// already prepared
	if _, ok := p.prepared[m.TableName()]; ok {
		return nil
	}

	p.Log.Info("** Preparing ** table >%s< statements", m.TableName())

	// get one
	query := m.GetOneSQL()

	getOneStmt, err := p.DB.Preparex(query)
	if err != nil {
		p.Log.Warn("Error preparing GetOneSQL statement >%v<", err)
		return err
	}

	getOneSQLList[m.TableName()] = query
	p.getOneStmtList[m.TableName()] = getOneStmt

	// get one for update
	query = m.GetOneForUpdateSQL()

	getOneForUpdateStmt, err := p.DB.Preparex(query)
	if err != nil {
		p.Log.Warn("Error preparing GetOneForUpdateSQL statement >%v<", err)
		return err
	}

	getOneForUpdateSQLList[m.TableName()] = query
	p.getOneForUpdateStmtList[m.TableName()] = getOneForUpdateStmt

	// get many
	query = m.GetManySQL()

	getManyStmt, err := p.DB.PrepareNamed(m.GetManySQL())
	if err != nil {
		p.Log.Warn("Error preparing GetManySQL statement >%v<", err)
		return err
	}

	getManySQLList[m.TableName()] = query
	p.getManyStmtList[m.TableName()] = getManyStmt

	// create
	query = m.CreateOneSQL()

	createStmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		p.Log.Warn("Error preparing CreateSQL statement >%v<", err)
		return err
	}

	createSQLList[m.TableName()] = query
	p.createStmtList[m.TableName()] = createStmt

	// update
	query = m.UpdateOneSQL()

	updateOneStmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		p.Log.Warn("Error preparing UpdateOneSQL statement >%v<", err)
		return err
	}

	updateOneSQLList[m.TableName()] = query
	p.updateOneStmtList[m.TableName()] = updateOneStmt

	// update many
	query = m.UpdateManySQL()

	updateManyStmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		p.Log.Warn("Error preparing UpdateManySQL statement >%v<", err)
		return err
	}

	updateManySQLList[m.TableName()] = query
	p.updateManyStmtList[m.TableName()] = updateManyStmt

	// delete
	query = m.DeleteOneSQL()

	deleteStmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		p.Log.Warn("Error preparing DeleteSQL statement >%v<", err)
		return err
	}

	deleteOneSQLList[m.TableName()] = query
	p.deleteOneStmtList[m.TableName()] = deleteStmt

	// delete many
	query = m.DeleteManySQL()

	deleteManyStmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		p.Log.Warn("Error preparing DeleteManySQL statement >%v<", err)
		return err
	}

	deleteManySQLList[m.TableName()] = query
	p.deleteManyStmtList[m.TableName()] = deleteManyStmt

	// remove
	query = m.RemoveOneSQL()

	removeStmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		p.Log.Warn("Error preparing RemoveSQL statement >%v<", err)
		return err
	}

	removeOneSQLList[m.TableName()] = query
	p.removeOneStmtList[m.TableName()] = removeStmt

	// remove many
	query = m.RemoveManySQL()

	removeManyStmt, err := p.DB.PrepareNamed(query)
	if err != nil {
		p.Log.Warn("Error preparing RemoveManySQL statement >%v<", err)
		return err
	}

	removeManySQLList[m.TableName()] = query
	p.removeManyStmtList[m.TableName()] = removeManyStmt

	p.prepared[m.TableName()] = true

	return nil
}

// GetOneStmt -
func (p *Prepare) GetOneStmt(m preparable.Preparable) *sqlx.Stmt {

	stmt := p.getOneStmtList[m.TableName()]

	return stmt
}

// GetOneForUpdateStmt -
func (p *Prepare) GetOneForUpdateStmt(m preparable.Preparable) *sqlx.Stmt {

	stmt := p.getOneForUpdateStmtList[m.TableName()]

	return stmt
}

// GetManyStmt -
func (p *Prepare) GetManyStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.getManyStmtList[m.TableName()]

	return stmt
}

// CreateOneStmt -
func (p *Prepare) CreateOneStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.createStmtList[m.TableName()]

	return stmt
}

// UpdateOneStmt -
func (p *Prepare) UpdateOneStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.updateOneStmtList[m.TableName()]

	return stmt
}

// UpdateManyStmt -
func (p *Prepare) UpdateManyStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.updateManyStmtList[m.TableName()]

	return stmt
}

// DeleteOneStmt -
func (p *Prepare) DeleteOneStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.deleteOneStmtList[m.TableName()]

	return stmt
}

// DeleteManyStmt -
func (p *Prepare) DeleteManyStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.deleteManyStmtList[m.TableName()]

	return stmt
}

// RemoveOneStmt -
func (p *Prepare) RemoveOneStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.removeOneStmtList[m.TableName()]

	return stmt
}

// RemoveManyStmt -
func (p *Prepare) RemoveManyStmt(m preparable.Preparable) *sqlx.NamedStmt {

	stmt := p.removeManyStmtList[m.TableName()]

	return stmt
}

// GetOneSQL -
func (p *Prepare) GetOneSQL(m preparable.Preparable) string {

	query := getOneSQLList[m.TableName()]

	return query
}

// GetManySQL -
func (p *Prepare) GetManySQL(m preparable.Preparable) string {

	query := getManySQLList[m.TableName()]

	return query
}

// CreateSQL -
func (p *Prepare) CreateSQL(m preparable.Preparable) string {

	query := createSQLList[m.TableName()]

	return query
}

// UpdateOneSQL -
func (p *Prepare) UpdateOneSQL(m preparable.Preparable) string {

	query := updateOneSQLList[m.TableName()]

	return query
}

// UpdateManySQL -
func (p *Prepare) UpdateManySQL(m preparable.Preparable) string {

	query := updateManySQLList[m.TableName()]

	return query
}

// DeleteOneSQL -
func (p *Prepare) DeleteOneSQL(m preparable.Preparable) string {

	query := deleteOneSQLList[m.TableName()]

	return query
}

// DeleteManySQL -
func (p *Prepare) DeleteManySQL(m preparable.Preparable) string {

	query := deleteManySQLList[m.TableName()]

	return query
}

// RemoveOneSQL -
func (p *Prepare) RemoveOneSQL(m preparable.Preparable) string {

	query := removeOneSQLList[m.TableName()]

	return query
}

// RemoveManySQL -
func (p *Prepare) RemoveManySQL(m preparable.Preparable) string {

	query := removeManySQLList[m.TableName()]

	return query
}
