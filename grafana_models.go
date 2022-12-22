package main

type Folders map[string][]RuleGroup

type RuleGroup struct {
	Name     string `json:"name"`
	Interval string `json:"interval"`
	Rules    []Rule `json:"rules"`
}

type Rule struct {
	Expr         string            `json:"expr"`
	For          string            `json:"for"`
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations"`
	GrafanaAlert GrafanaAlert      `json:"grafana_alert"`
}

type GrafanaAlert struct {
	ID              int64   `json:"id"`
	OrgID           int64   `json:"orgId"`
	Title           string  `json:"title"`
	Condition       string  `json:"condition"`
	Data            []Datum `json:"data"`
	Updated         string  `json:"updated"`
	IntervalSeconds int64   `json:"intervalSeconds"`
	Version         int64   `json:"version"`
	Uid             string  `json:"uid"`
	NamespaceUid    string  `json:"namespace_uid"`
	NamespaceID     int64   `json:"namespace_id"`
	RuleGroup       string  `json:"rule_group"`
	NoDataState     string  `json:"no_data_state"`
	ExecErrState    string  `json:"exec_err_state"`
	Provenance      string  `json:"provenance"`
}

type Datum struct {
	RefID             string            `json:"refId"`
	QueryType         string            `json:"queryType"`
	RelativeTimeRange RelativeTimeRange `json:"relativeTimeRange"`
	DatasourceUid     string            `json:"datasourceUid"`
	Model             map[string]any    `json:"model"`
}

type Condition struct {
	Evaluator Evaluator `json:"evaluator"`
	Operator  Operator  `json:"operator"`
	Query     Query     `json:"query"`
	Reducer   Evaluator `json:"reducer"`
	Type      string    `json:"type"`
}

type Evaluator struct {
	Params []int64 `json:"params"`
	Type   string  `json:"type"`
}

type Operator struct {
	Type string `json:"type"`
}

type Query struct {
	Params []string `json:"params"`
}

type Datasource struct {
	ID   *int64  `json:"id,omitempty"`
	Type string  `json:"type"`
	Uid  *string `json:"uid,omitempty"`
}

type Dimensions struct {
	QueueName string `json:"QueueName"`
}

type RelativeTimeRange struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}
