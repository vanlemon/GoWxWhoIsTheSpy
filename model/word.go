package model

import (
	"context"
	"sync"

	"lmf.mortal.com/GoLogs"
)

// 词汇信息模型
type Word struct {
	ID         int    `gorm:"column:id"`          // 自增ID
	NormalWord string `gorm:"column:normal_word"` // 平民词汇，平民可以没有词（白板），即只有卧底有词
	SpyWord    string `gorm:"column:spy_word"`    // 卧底词汇
	BlankWord  string `gorm:"column:blank_word"`  // 空白词汇，白板一定没有词，可以是类别提示
	Class      string `gorm:"column:class"`       // 词汇类别
}

// 词汇信息表名
func (Word) TableName() string {
	return "word"
}

// 词汇模型访问对象
type WordDao struct {
}

// 词汇模型访问对象 - 单例模式
var wordDao *WordDao
var wordDaoOnce sync.Once

// 词汇模型访问示例
func WordDaoInstance() *WordDao {
	wordDaoOnce.Do(func() {
		wordDao = &WordDao{}
	})
	return wordDao
}

// 创建词汇
func (d *WordDao) CreateWord(ctx context.Context, word *Word) (err error) {
	logs.CtxInfo(ctx, "[Model MySQL CreateWord Req] req: %#v", word) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model MySQL CreateWord Resp] resp: %#v", err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	err = WordDB.Model(&Word{}).Create(word).Error

	if err != nil {
		logs.CtxError(ctx, "[Model MySQL CreateWord] Create error: %#v", err)
		return err
	}
	/************************************ 核心逻辑 ****************************************/

	return nil
}

// 随机获取词汇 - 一条
func (d *WordDao) RandomQueryWord(ctx context.Context) (word *Word, err error) {
	logs.CtxInfo(ctx, "[Model MySQL RandomQueryWord Req] req") // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model MySQL RandomQueryWord Resp] resp: %#v, %#v", word, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	word = &Word{} // 定义接收结构体
	// 随机查询一个词汇
	err = WordDB.Model(&Word{}).Order("RAND()").Take(word).Error

	if err != nil {
		logs.CtxError(ctx, "[Model MySQL RandomQueryWord] Query error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return word, nil
}

// 随机获取词汇 - 多条
func (d *WordDao) RandomQueryWordList(ctx context.Context, count int) (wordList []*Word, err error) {
	logs.CtxInfo(ctx, "[Model MySQL RandomQueryWordList Req] req: %#v", count) // 入口日志
	defer func() {
		logs.CtxInfo(ctx, "[Model MySQL RandomQueryWordList Resp] resp: %#v, %#v", wordList, err) // 出口日志
	}()

	/************************************ 核心逻辑 ****************************************/
	wordList = []*Word{} // 定义接收结构体
	// 随机查询一个词汇
	err = WordDB.Model(&Word{}).Order("RAND()").Limit(count).Find(&wordList).Error

	if err != nil {
		logs.CtxError(ctx, "[Model MySQL RandomQueryWordList] Query error: %#v", err)
		return nil, err
	}
	/************************************ 核心逻辑 ****************************************/

	return wordList, nil
}
