package app

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NgeKaworu/time-mgt-go/src/models"
	"github.com/NgeKaworu/time-mgt-go/src/parsup"
	"github.com/NgeKaworu/time-mgt-go/src/resultor"
	"github.com/NgeKaworu/time-mgt-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddTag 添加标签
func (d *App) AddTag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	if len(body) == 0 {
		resultor.RetFail(w, errors.New("not has body"))
		return
	}

	p, err := parsup.ParSup().ConvJSON(body)
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	err = utils.Required(p, map[string]string{
		"name":  "标签名不能为空",
		"color": "颜色不能为空",
	})

	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	t := d.mongo.GetColl(models.TTag)
	p["uid"] = uid
	p["createAt"] = time.Now().Local()

	res, err := t.InsertOne(context.Background(), p)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "dup key") {
			errMsg = "该标签已被创建"
		}

		resultor.RetFail(w, errors.New(errMsg))
		return
	}

	resultor.RetOk(w, res.InsertedID.(primitive.ObjectID).Hex())
}

// SetTag 更新标签
func (d *App) SetTag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	if len(body) == 0 {
		resultor.RetFail(w, errors.New("not has body"))
		return
	}

	p, err := parsup.ParSup().ConvJSON(body)
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	err = utils.Required(p, map[string]string{
		"id":    "标签不能id为空",
		"name":  "标签名不能为空",
		"color": "颜色不能为空",
	})
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	t := d.mongo.GetColl(models.TTag)
	p["uid"] = uid
	p["updateAt"] = time.Now().Local()

	id := p["id"]
	delete(p, "id")

	res := t.FindOneAndUpdate(context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": p},
	)
	if res.Err() != nil {
		resultor.RetFail(w, res.Err())
		return
	}

	resultor.RetOk(w, "修改成功")
}

// RemoveTag 删除标签
func (d *App) RemoveTag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	t := d.mongo.GetColl(models.TRecord)

	used, err := t.CountDocuments(context.Background(), bson.M{
		"uid": uid,
		"tid": id,
	})

	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	if used != 0 {
		resultor.RetFail(w, errors.New("不能删除正在使用的标签。"))
		return
	}

	t = d.mongo.GetColl(models.TTag)

	res := t.FindOneAndDelete(context.Background(), bson.M{"_id": id, "uid": uid})

	if res.Err() != nil {
		resultor.RetFail(w, res.Err())
		return
	}

	resultor.RetOk(w, "删除成功")
}

// ListTag tag列表
func (d *App) ListTag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := r.URL.Query()
	l := q.Get("limit")
	s := q.Get("skip")

	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	limit, _ := strconv.ParseInt(l, 10, 64)
	skip, _ := strconv.ParseInt(s, 10, 64)

	t := d.mongo.GetColl(models.TTag)

	cur, err := t.Find(context.Background(), bson.M{
		"uid": uid,
	}, options.Find().SetSkip(skip).SetLimit(limit))

	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	list := make([]models.Tag, 0)

	err = cur.All(context.Background(), &list)
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	resultor.RetOk(w, list)
}
