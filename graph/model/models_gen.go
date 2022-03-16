// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ClusterCockpit/cc-backend/schema"
)

type Accelerator struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Model string `json:"model"`
}

type Count struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type FilterRanges struct {
	Duration  *IntRangeOutput  `json:"duration"`
	NumNodes  *IntRangeOutput  `json:"numNodes"`
	StartTime *TimeRangeOutput `json:"startTime"`
}

type FloatRange struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

type Footprints struct {
	Nodehours []schema.Float      `json:"nodehours"`
	Metrics   []*MetricFootprints `json:"metrics"`
}

type HistoPoint struct {
	Count int `json:"count"`
	Value int `json:"value"`
}

type IntRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type IntRangeOutput struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type JobFilter struct {
	Tags            []string          `json:"tags"`
	JobID           *StringInput      `json:"jobId"`
	ArrayJobID      *int              `json:"arrayJobId"`
	User            *StringInput      `json:"user"`
	Project         *StringInput      `json:"project"`
	Cluster         *StringInput      `json:"cluster"`
	Partition       *StringInput      `json:"partition"`
	Duration        *IntRange         `json:"duration"`
	MinRunningFor   *int              `json:"minRunningFor"`
	NumNodes        *IntRange         `json:"numNodes"`
	NumAccelerators *IntRange         `json:"numAccelerators"`
	NumHWThreads    *IntRange         `json:"numHWThreads"`
	StartTime       *TimeRange        `json:"startTime"`
	State           []schema.JobState `json:"state"`
	FlopsAnyAvg     *FloatRange       `json:"flopsAnyAvg"`
	MemBwAvg        *FloatRange       `json:"memBwAvg"`
	LoadAvg         *FloatRange       `json:"loadAvg"`
	MemUsedMax      *FloatRange       `json:"memUsedMax"`
}

type JobMetricWithName struct {
	Name   string            `json:"name"`
	Metric *schema.JobMetric `json:"metric"`
}

type JobResultList struct {
	Items  []*schema.Job `json:"items"`
	Offset *int          `json:"offset"`
	Limit  *int          `json:"limit"`
	Count  *int          `json:"count"`
}

type JobsStatistics struct {
	ID             string        `json:"id"`
	TotalJobs      int           `json:"totalJobs"`
	ShortJobs      int           `json:"shortJobs"`
	TotalWalltime  int           `json:"totalWalltime"`
	TotalCoreHours int           `json:"totalCoreHours"`
	HistWalltime   []*HistoPoint `json:"histWalltime"`
	HistNumNodes   []*HistoPoint `json:"histNumNodes"`
}

type MetricConfig struct {
	Name     string             `json:"name"`
	Unit     string             `json:"unit"`
	Scope    schema.MetricScope `json:"scope"`
	Timestep int                `json:"timestep"`
	Peak     float64            `json:"peak"`
	Normal   float64            `json:"normal"`
	Caution  float64            `json:"caution"`
	Alert    float64            `json:"alert"`
}

type MetricFootprints struct {
	Metric string         `json:"metric"`
	Data   []schema.Float `json:"data"`
}

type NodeMetrics struct {
	Host    string               `json:"host"`
	Metrics []*JobMetricWithName `json:"metrics"`
}

type OrderByInput struct {
	Field string            `json:"field"`
	Order SortDirectionEnum `json:"order"`
}

type PageRequest struct {
	ItemsPerPage int `json:"itemsPerPage"`
	Page         int `json:"page"`
}

type StringInput struct {
	Eq         *string `json:"eq"`
	Contains   *string `json:"contains"`
	StartsWith *string `json:"startsWith"`
	EndsWith   *string `json:"endsWith"`
}

type SubCluster struct {
	Name            string    `json:"name"`
	Nodes           string    `json:"nodes"`
	ProcessorType   string    `json:"processorType"`
	SocketsPerNode  int       `json:"socketsPerNode"`
	CoresPerSocket  int       `json:"coresPerSocket"`
	ThreadsPerCore  int       `json:"threadsPerCore"`
	FlopRateScalar  int       `json:"flopRateScalar"`
	FlopRateSimd    int       `json:"flopRateSimd"`
	MemoryBandwidth int       `json:"memoryBandwidth"`
	Topology        *Topology `json:"topology"`
}

type TimeRange struct {
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
}

type TimeRangeOutput struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Topology struct {
	Node         []int          `json:"node"`
	Socket       [][]int        `json:"socket"`
	MemoryDomain [][]int        `json:"memoryDomain"`
	Die          [][]int        `json:"die"`
	Core         [][]int        `json:"core"`
	Accelerators []*Accelerator `json:"accelerators"`
}

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type Aggregate string

const (
	AggregateUser    Aggregate = "USER"
	AggregateProject Aggregate = "PROJECT"
	AggregateCluster Aggregate = "CLUSTER"
)

var AllAggregate = []Aggregate{
	AggregateUser,
	AggregateProject,
	AggregateCluster,
}

func (e Aggregate) IsValid() bool {
	switch e {
	case AggregateUser, AggregateProject, AggregateCluster:
		return true
	}
	return false
}

func (e Aggregate) String() string {
	return string(e)
}

func (e *Aggregate) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Aggregate(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Aggregate", str)
	}
	return nil
}

func (e Aggregate) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortDirectionEnum string

const (
	SortDirectionEnumDesc SortDirectionEnum = "DESC"
	SortDirectionEnumAsc  SortDirectionEnum = "ASC"
)

var AllSortDirectionEnum = []SortDirectionEnum{
	SortDirectionEnumDesc,
	SortDirectionEnumAsc,
}

func (e SortDirectionEnum) IsValid() bool {
	switch e {
	case SortDirectionEnumDesc, SortDirectionEnumAsc:
		return true
	}
	return false
}

func (e SortDirectionEnum) String() string {
	return string(e)
}

func (e *SortDirectionEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortDirectionEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortDirectionEnum", str)
	}
	return nil
}

func (e SortDirectionEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
