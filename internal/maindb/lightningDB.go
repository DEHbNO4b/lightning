package maindb

import (
	"database/sql"
	"strconv"

	"github.com/DEHbNO4b/lightning.git/internal/domain/models"
)

var dsn string = "postgres://postgres:917836@localhost:5432/lightning?"

var queryMakeTab string = ` DROP TABLE IF EXISTS strikes;
					CREATE TABLE IF NOT EXISTS strikes(
						id serial primary key,
						time timestamptz,
						latitude numeric(6,4),
						longitude numeric(6,4),
						geog GEOGRAPHY(Point),
						signal smallint,
						cloud boolean,
						cluster integer
					);
					CREATE INDEX ON strikes USING GIST(geog);`

var queryInsert string = `INSERT INTO strikes (time,longitude,latitude,geog,signal,cloud) 
							VALUES($1,$2,$3,ST_MakePoint($4, $5)::GEOGRAPHY,$6,$7)
							RETURNING ID;`
var queryNeighbors string = `SELECT id,longitude,latitude FROM strikes WHERE geog<->st_setSRID(st_makePoint($1,$2),4326)::GEOGRAPHY < $3;`

type LightningDB struct {
	Db *sql.DB
}

//	func NewLightningDB(db *sql.DB) LightningDB {
//		return LightningDB{Db: db}
//	}
func (n *LightningDB) GetNeighbourse(long, lat float32, eps int) (map[string]models.Stroke, error) {
	var ans = make(map[string]models.Stroke)
	var latitude, longitude float32
	var id int
	rows, err := n.Db.Query(queryNeighbors, long, lat, eps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &longitude, &latitude); err != nil {
			return nil, err
		}
		s := models.Stroke{Id: id, Longitude: longitude, Latitude: latitude}
		ans[strconv.Itoa(id)] = s
	}
	return ans, nil
}

func (n *LightningDB) OpenDB() error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	n.Db = db
	return nil
}
func (n *LightningDB) MakeTab() error {
	_, err := n.Db.Exec(queryMakeTab)
	if err != nil {
		return err
	}
	return nil
}
func (n *LightningDB) LoadRawToDb(strokes []models.Stroke) (map[string]models.Stroke, error) {
	data := make(map[string]models.Stroke, len(strokes))
	for i, el := range strokes {
		var idInDB int
		err := n.Db.QueryRow(queryInsert, el.Time, el.Longitude, el.Latitude, el.Longitude, el.Latitude, el.Signal, el.Cloud).Scan(&idInDB)
		if err != nil {
			return nil, err
		}
		strokes[i].Id = idInDB
		data[strconv.Itoa(idInDB)] = strokes[i]
	}

	return data, nil
}
func (n *LightningDB) LoadClasterToDb(data map[string]models.Stroke) error {
	for key, el := range data {
		id, _ := strconv.Atoi(key)
		_, err := n.Db.Exec(`UPDATE strikes SET cluster = $1 WHERE id = $2`, el.Claster, id)
		if err != nil {
			return err
		}
	}
	return nil
}
