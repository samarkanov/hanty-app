package model

import (
    "sync"
    // "fmt"
    "encoding/json"
)

type db_row struct {
    Topic string
    Data []string
}

type db_rows = []db_row

type db_table struct {
    Client string
    Data []db_row
}

type DB []db_table

type Query struct {
    Client string
    Topic string
    Entry string
}

/* Singleton */
var (
    db * DB
    once sync.Once
)

func get_db() * DB {
    once.Do(func() {
        db = init_db()
    })

    return db
}

/* Config Constructor */
func Cursor() * DB {
    return get_db()
}

func init_db() * DB {
    var res DB
    return &res
}

func (DB) remove(idx int) {
    dbb := *db
    dbb[len(dbb)-1], dbb[idx] = dbb[idx], dbb[len(dbb)-1]
    *db = dbb[:len(dbb)-1]
}

func (DB) remove_row(idx int, jdx int) {
    dbb := (*db)[idx].Data
    len_ := len(dbb) - 1
    dbb[len_], dbb[jdx] = dbb[jdx], dbb[len_]
    (*db)[idx].Data = dbb[:len_]
}

func (DB) remove_entry(idx int, jdx int, edx int) {
    dbb := (*db)[idx].Data[jdx].Data
    len_ := len(dbb) - 1
    dbb[len_], dbb[edx] = dbb[edx], dbb[len_]
    (*db)[idx].Data[jdx].Data = dbb[:len_]
}


/* Public interfaces */
func (DB) Set(query * Query) {
    // set and update
    var table_idx = -1

    for idx, table := range *db {
        if table.Client == query.Client {
            table_idx = idx
            for jdx, row := range table.Data {
                if row.Topic == query.Topic {
                    // do not insert if entry already exists:
                    for _, entry := range row.Data {
                        if entry == query.Entry {
                            return
                        }
                    }
                    ptr := (*db)[idx].Data[jdx].Data
                    (*db)[idx].Data[jdx].Data = append(ptr, query.Entry)
                    return
                }
            }
        }
    }

    row := db_row{
        Topic: query.Topic,
        Data: []string{query.Entry},
    }

    if table_idx == -1 {
        // create table
        table := db_table{
            Client: query.Client,
            Data: []db_row{row},
        }
        *db = append(*db, table)
    } else {
        // create row
        ptr := (*db)[table_idx].Data
        (*db)[table_idx].Data = append(ptr, row)
    }
}

func (DB) Delete(query * Query) {
    if len(query.Topic) > 0 {
        // remove row or entry in a row
        for idx, table := range *db {
            if table.Client == query.Client {
                for jdx, row := range table.Data {
                    if row.Topic == query.Topic {
                        if len(query.Entry) > 0 {
                            // remove entry in a row
                            for edx, entry_ := range row.Data{
                                if query.Entry == entry_ {
                                    db.remove_entry(idx, jdx, edx)
                                    return
                                }
                            }
                        } else {
                            // remove row
                            db.remove_row(idx, jdx)
                            return
                        }
                    }
                }
            }
        }
    } else {
        // remove table
        for idx, table := range *db {
            if table.Client == query.Client {
                db.remove(idx)
                return
            }
        }
    }
}

func empty() (string) {
    var res_ map[string]string
    res, _ := json.Marshal(res_)
    return string(res)
}

func (DB) Get(query * Query) (string) {
    // if topic is part of the query parameters:
    if len(query.Topic) > 0 {
        // getting row
        for _, table := range *db {
            if table.Client == query.Client {
                for _, row := range table.Data {
                    if row.Topic == query.Topic {
                        // If the client and the topic are matching
                        // with the DB entries: return it
                        res, _ := json.Marshal(row.Data)
                        return string(res)
                    }
                }
            }
        }
    } else {
        // getting table
        for _, table := range *db {
            if table.Client == query.Client {
                if len(table.Data) == 0 {
                    // no entries in the DB: return empty json (not null)
                    return empty()
                }
                res, _ := json.Marshal(table.Data)
                return string(res)
            }
        }
    }

    return empty()
}
