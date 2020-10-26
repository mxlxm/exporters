package rocketmq

import (
	"strconv"
	"sync"

	"github.com/mxlxm/RocketmqExporter/model"
	"github.com/mxlxm/RocketmqExporter/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const metricsPrefix = "rocketmq"

type Options struct {
	Registry *prometheus.Registry
}

var (
	log *zap.SugaredLogger
)

func SetLogger(l *zap.SugaredLogger) {
	log = l
}

func DefaultOptions() Options {
	return Options{}
}

type Exporter struct {
	target                       string
	msgDiffDetail                prometheus.GaugeVec
	msgDiffTopic                 prometheus.GaugeVec
	msgDiffConsumerGroup         prometheus.GaugeVec
	msgDiffTopicAndConsumerGroup prometheus.GaugeVec
	msgDiffBroker                prometheus.GaugeVec
	msgDiffQueue                 prometheus.GaugeVec
	msgDiffClientinfo            prometheus.GaugeVec
}

func NewRocketMQExporter(target string, opts Options) error {

	msgDiffDetail := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "msg_diff_detail",
		Help:      "msg diff detail group by every consumer client",
	}, []string{"topic", "consumerGroup", "broker", "queueId", "consumerClientIP", "consumerClientPID"})

	msgDiffTopic := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "msg_diff_topic",
		Help:      "msg diff group by every topic",
	}, []string{"topic"})

	msgDiffConsumerGroup := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "msg_diff_consumergroup",
		Help:      "msg diff group by every consumer group",
	}, []string{"consumerGroup"})

	msgDiffTopicAndConsumerGroup := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "msg_diff_topic_consumergroup",
		Help:      "msg diff group by topic and consumer group",
	}, []string{"topic", "consumerGroup"})

	msgDiffBroker := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "msg_diff_broker",
		Help:      "msg diff group by broker",
	}, []string{"broker"})

	msgDiffQueue := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "msg_diff_queue",
		Help:      "msg diff group by broker and queue",
	}, []string{"broker", "queueId"})

	msgDiffClientinfo := *prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsPrefix,
		Name:      "msg_diff_clientinfo",
		Help:      "msg diff group by clientinfo",
	}, []string{"consumerClientIP", "consumerClientPID"})

	e := &Exporter{
		target:                       target,
		msgDiffDetail:                msgDiffDetail,
		msgDiffTopic:                 msgDiffTopic,
		msgDiffConsumerGroup:         msgDiffConsumerGroup,
		msgDiffTopicAndConsumerGroup: msgDiffTopicAndConsumerGroup,
		msgDiffBroker:                msgDiffBroker,
		msgDiffQueue:                 msgDiffQueue,
		msgDiffClientinfo:            msgDiffClientinfo,
	}
	opts.Registry.MustRegister(e)
	return nil
}

func (exporter *Exporter) Collect(ch chan<- prometheus.Metric) {

	var msgDiff *model.MsgDiff = service.MsgUnconsumedCount(exporter.target)

	if msgDiff != nil && msgDiff.MsgDiff_Details != nil {
		var msg_Diff_Detail_Array []*model.MsgDiff_Detail = msgDiff.MsgDiff_Details
		for _, e := range msg_Diff_Detail_Array {
			exporter.msgDiffDetail.WithLabelValues(e.Topic, e.ConsumerGroup, e.Broker, strconv.Itoa(e.QueueId), e.ConsumerClientIP, e.ConsumerClientPID).Set(float64(e.Diff))
		}
	}

	if msgDiff != nil && msgDiff.MsgDiff_Topics != nil {
		var msg_Diff_Topic_Map map[string]*model.MsgDiff_Topic = msgDiff.MsgDiff_Topics
		for _, v := range msg_Diff_Topic_Map {
			exporter.msgDiffTopic.WithLabelValues(v.Topic).Set(float64(v.Diff))
		}
	}

	if msgDiff != nil && msgDiff.MsgDiff_ConsumerGroups != nil {
		var msg_Diff_ConsumerGroup_Map map[string]*model.MsgDiff_ConsumerGroup = msgDiff.MsgDiff_ConsumerGroups
		for _, v := range msg_Diff_ConsumerGroup_Map {
			exporter.msgDiffConsumerGroup.WithLabelValues(v.ConsumerGroup).Set(float64(v.Diff))
		}
	}

	if msgDiff != nil && msgDiff.MsgDiff_Topics_ConsumerGroups != nil {
		var msg_Diff_Topic_ConsumerGroup_Map map[string]*model.MsgDiff_Topic_ConsumerGroup = msgDiff.MsgDiff_Topics_ConsumerGroups
		for _, v := range msg_Diff_Topic_ConsumerGroup_Map {
			exporter.msgDiffTopicAndConsumerGroup.WithLabelValues(v.Topic, v.ConsumerGroup).Set(float64(v.Diff))
		}
	}

	if msgDiff != nil && msgDiff.MsgDiff_Brokers != nil {
		var msg_Diff_Broker_Map map[string]*model.MsgDiff_Broker = msgDiff.MsgDiff_Brokers
		for _, v := range msg_Diff_Broker_Map {
			exporter.msgDiffBroker.WithLabelValues(v.Broker).Set(float64(v.Diff))
		}
	}

	if msgDiff != nil && msgDiff.MsgDiff_Queues != nil {
		var msg_Diff_Queue_Map map[string]*model.MsgDiff_Queue = msgDiff.MsgDiff_Queues
		for _, v := range msg_Diff_Queue_Map {
			exporter.msgDiffQueue.WithLabelValues(v.Broker, strconv.Itoa(v.QueueId)).Set(float64(v.Diff))
		}
	}

	if msgDiff != nil && msgDiff.MsgDiff_ClientInfos != nil {
		var msg_Diff_Clientinfo_Map map[string]*model.MsgDiff_ClientInfo = msgDiff.MsgDiff_ClientInfos
		for _, v := range msg_Diff_Clientinfo_Map {
			exporter.msgDiffClientinfo.WithLabelValues(v.ConsumerClientIP, v.ConsumerClientPID)
		}
	}

	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		defer wg.Done()
		exporter.msgDiffDetail.Collect(ch)
	}()
	go func() {
		defer wg.Done()
		exporter.msgDiffTopic.Collect(ch)
	}()
	go func() {
		defer wg.Done()
		exporter.msgDiffConsumerGroup.Collect(ch)
	}()
	go func() {
		defer wg.Done()
		exporter.msgDiffTopicAndConsumerGroup.Collect(ch)
	}()
	go func() {
		defer wg.Done()
		exporter.msgDiffBroker.Collect(ch)
	}()
	go func() {
		defer wg.Done()
		exporter.msgDiffQueue.Collect(ch)
	}()
	go func() {
		defer wg.Done()
		exporter.msgDiffClientinfo.Collect(ch)
	}()
	wg.Wait()
}

func (exporter *Exporter) Describe(ch chan<- *prometheus.Desc) {
	exporter.msgDiffDetail.Describe(ch)
}
