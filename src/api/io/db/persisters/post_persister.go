package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"fmt"
	"log"
)

type PostPersister struct {
	connection *db.Connection
}

func NewPostPersister(connection *db.Connection) *PostPersister {
	return &PostPersister{connection}
}

func (p *PostPersister) Insert(post *model.Post) (*model.Post, error) {
	op := Operations{connection: p.connection}
	postModel, err := op.Insert(post.Url, post.Title, post.Location)
	if err != nil {
		return nil, err
	}
	return parseToPost(postModel), nil
}

func (p *PostPersister) Get(id int64) (*model.Post, error) {
	op := Operations{connection: p.connection}
	query := "SELECT * FROM tpms_prod.posts WHERE id = ?"
	postModel, err := op.Get(query, id)
	if err != nil {
		return nil, err
	}
	return parseToPost(&postModel[0]), nil
}

func (p *PostPersister) GetAll() ([]model.Post, error) {
	op := Operations{connection: p.connection}
	resultList, err := op.Get("SELECT * FROM tpms_prod.posts")
	if err != nil {
		return nil, err
	}
	return parseToPostList(resultList), nil
}

func (p *PostPersister) Update(post *model.Post) (*model.Post, error) {
	op := Operations{connection: p.connection}
	query := "UPDATE tpms_prod.posts SET url = ?, title = ?, location = ? WHERE id = ?"
	postModel, err := op.Update(query, post.Url, post.Title, post.Location, post.Id)
	if err != nil {
		return nil, err
	}
	return parseToPost(postModel), nil
}

func (p *PostPersister) Delete(id int64) (bool, error) {
	op := Operations{connection: p.connection}
	deleted := op.Delete("DELETE FROM tpms_prod.posts WHERE id = ?", id)
	if !deleted {
		return false, fmt.Errorf("error deleting post with id %d", id)
	}
	return true, nil
}

func (p *PostPersister) DeleteByDogId(id int64) (bool, error) {
	op := Operations{connection: p.connection}
	deleted := op.Delete("DELETE FROM tpms_prod.posts WHERE dog_id = ?", id)
	if !deleted {
		return false, fmt.Errorf("error deleting post with dogID %d", id)
	}
	return true, nil
}

type Operations struct {
	connection *db.Connection
}

func (o *Operations) Insert(params ...any) (*model.PostModel, error) {
	result, err := o.connection.DB.Exec("INSERT INTO tpms_prod.posts(url, title, location) VALUES (?, ?, ?)", params)
	if err != nil {
		log.Printf("error insterting post: %s", err.Error())
		return nil, err
	}
	id, _ := result.LastInsertId()
	posts, _ := o.Get("select * from posts where id = ?", id)
	return &posts[0], nil
}

func (o *Operations) Get(stm string, params ...any) ([]model.PostModel, error) {
	rows, err := o.connection.DB.Query(stm, params)
	if err != nil {
		log.Printf("error getting posts: %s", err.Error())
		return nil, err
	}

	var resultList []model.PostModel
	for rows.Next() {
		var postModel model.PostModel
		if err = rows.Scan(&postModel.Id, &postModel.DogId, &postModel.Url, &postModel.Title, &postModel.Location); err != nil {
			log.Printf("error parsing post: %s", err.Error())
			return nil, err
		}
	}
	return resultList, nil
}

func (o *Operations) Update(stm string, params ...any) (*model.PostModel, error) {
	result, err := o.connection.DB.Exec(stm, params)
	if err != nil {
		log.Printf("error insterting post: %s", err.Error())
		return nil, err
	}
	id, _ := result.LastInsertId()
	posts, _ := o.Get("selecct * from posts where id = ?", id)
	return &posts[0], nil
}

func (o *Operations) Delete(stm string, params ...any) bool {
	_, err := o.connection.DB.Exec(stm, params)
	if err != nil {
		log.Printf("error deleting post: %s", err.Error())
		return false
	}
	return true
}

func parseToPost(postModel *model.PostModel) *model.Post {
	id, dogID := int64(0), int64(0)
	if postModel.Id.Valid {
		id = postModel.Id.Int64
	}
	if postModel.Id.Valid {
		id = postModel.Id.Int64
	}

	return &model.Post{
		Id:       id,
		DogId:    dogID,
		Title:    postModel.Title,
		Url:      postModel.Url,
		Location: postModel.Location,
	}
}

func parseToPostList(list []model.PostModel) []model.Post {
	if list == nil || len(list) == 0 {
		return make([]model.Post, 0)
	}
	var resultList []model.Post
	for _, l := range list {
		resultList = append(resultList, *parseToPost(&l))
	}
	return resultList
}
