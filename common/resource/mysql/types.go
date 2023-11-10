package mysql

import "strconv"

/**
基于PDM,建立统计的数据domain结构。注意，不要选择0， 因为go的int默认值为0
*/

//主键ID
type ID = uint64

//性别
type Sex = uint8

//加减
type AddReduce = uint8

//拥有
type YesNo = uint8

//短描述
type Str = string

//时间戳
type Time = int64

//状态
type UserYesNo = uint8

//平台
type Platform = uint8

//多枚举类型
type Type = uint8

//数量
type Num = uint32

//时间戳
type Timestamp = uint64

//排序
type Index = uint16

//数量，并且用到-1作为特殊标记位
type NumAll = int

//开启关闭
type OpenClose = uint8

//逻辑删除
type LogicDel = uint8

//设备
type Device = uint8

type PeriodType = uint8

type FinishYesNo = uint8

//性别
const (
	MAN   Sex = 1
	WOMAN Sex = 2
	EMPTY Sex = 0
)

//yes no
const (
	YES YesNo = 1
	NO  YesNo = 2
)

const (
	OPEN  OpenClose = 1
	CLOSE OpenClose = 2
)

//加，减
const (
	ADD    AddReduce = 1
	REDUCE AddReduce = 2
	SET    AddReduce = 3 // set
)

const (
	USER   UserYesNo = 1
	NOUSER UserYesNo = 2
)

func TypeToString(t Type) string {
	return strconv.Itoa(int(t))
}

//逻辑删除
const (
	EXIST LogicDel = 1
	DEL   LogicDel = 2
)

// 礼物类型
const (
	DiamondYellow Type = 1 // 黄钻礼物
	DiamondPink   Type = 2 // 粉钻礼物
)
