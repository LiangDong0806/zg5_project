package api

//
//func Preheat(c *gin.Context) { //TODO 预热数据
//	preheat, _ := utils.PreheatTheProduct()
//	for _, g := range *preheat {
//		for i := 0; i < g.GoodsStock; i++ {
//			key := "goods" + strconv.Itoa(int(g.GoodsID)) + ":" + g.GoodsDescription
//			err := utils.RdbLPush(key, i)
//			if err != nil {
//				c.JSON(404, gin.H{"error": "redis预热失败" + err.Error()})
//			}
//		}
//		fmt.Println(strconv.Itoa(g.GoodsStock), "][]]]]]]]]]]]]]]]")
//	}
//	c.JSON(http.StatusOK, "预热成功")
//
//}
//
//type Order struct {
//	ID        string
//	ProductID string
//}
//
//func SecJi(c *gin.Context) {
//	goodsId := c.PostForm("goods_id")
//	goodsDescription := c.PostForm("goodsDescription")
//	uid := c.PostForm("uid")
//	gid, _ := strconv.Atoi(goodsId)
//	key := "goods" + goodsId + ":" + goodsDescription
//
//	goodsStock, err := utils.RdbLLen(key)
//	if err != nil {
//		c.JSONP(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "Redis LLen取值失败" + err.Error()})
//		return
//	}
//	if goodsStock <= 0 {
//		c.JSONP(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "没库存了! 秒杀即将结束" + err.Error()})
//		return
//	}
//	re, err := utils.RdbLPop(key)
//	if err != nil {
//		c.JSONP(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "库存扣减失败！" + err.Error()})
//		return
//	}
//	fmt.Println(re, "redis:---")
//	// TODO: 下面开始扣减数据库库存
//	res, err := utils.CheckInventory(gid)
//	if err != nil {
//		c.JSONP(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "商品信息获取失败"})
//		return
//	}
//	res.GoodsStock--
//	err = utils.DeductionOfInventory(res)
//	if err != nil {
//		c.JSONP(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "库存扣减失败"})
//	}
//	order := Order{
//		ID:        uid,
//		ProductID: goodsId,
//	}
//
//	go func() {
//		rabbit := NewRabbitMQSimple("order")
//		rabbit.PublishSimple(order)
//	}()
//	data := map[string]interface{}{
//		"code":  200,
//		"msg":   "订单创建成功",
//		"order": order,
//	}
//	c.JSON(http.StatusOK, data)
//}
//
//// 连接信息amqp://kuteng:kuteng@127.0.0.1:5672/kuteng这个信息是固定不变的amqp://事固定参数后面两个是用户名密码ip地址端口号Virtual Host
//const MQURL = "amqp://guest:guest@127.0.0.1:5672"
//
//// rabbitMQ结构体
//type RabbitMQ struct {
//	conn    *amqp.Connection
//	channel *amqp.Channel
//	//队列名称
//	QueueName string
//	//交换机名称
//	Exchange string
//	//bind Key 名称
//	Key string
//	//连接信息
//	Mqurl string
//}
//
//// 创建结构体实例
//func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
//	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
//}
//
//// 错误处理函数
//func (r *RabbitMQ) failOnErr(err error, message string) {
//	if err != nil {
//		log.Fatalf("%s:%s", message, err)
//	}
//}
//
//// 创建简单模式下RabbitMQ实例
//func NewRabbitMQSimple(queueName string) *RabbitMQ {
//	//创建RabbitMQ实例
//	rabbitmq := NewRabbitMQ(queueName, "", "")
//	var err error
//	//获取connection
//	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
//	rabbitmq.failOnErr(err, "failed to connect rabb"+
//		"itmq!")
//	//获取channel
//	rabbitmq.channel, err = rabbitmq.conn.Channel()
//	rabbitmq.failOnErr(err, "failed to open a channel")
//	return rabbitmq
//}
//
//// 直接模式队列生产
//func (r *RabbitMQ) PublishSimple(message Order) {
//	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
//	_, err := r.channel.QueueDeclare(
//		r.QueueName,
//		//是否持久化
//		false,
//		//是否自动删除
//		false,
//		//是否具有排他性
//		false,
//		//是否阻塞处理
//		false,
//		//额外的属性
//		nil,
//	)
//	if err != nil {
//		fmt.Println(err)
//	}
//	buf := new(bytes.Buffer)
//	binary.Write(buf, binary.LittleEndian, &message)
//	bytes := buf.Bytes()
//	//调用channel 发送消息到队列中
//	r.channel.Publish(
//		r.Exchange,
//		r.QueueName,
//		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
//		false,
//		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
//		false,
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        bytes,
//		})
//}
//
//func CheckRegularly(cs *gin.Context) { //检查数据进行回滚
//	c := cron.New()
//
//	key1 := "goods" + "112" + ":" + "吉普冲锋衣" //TODO Redis里的第一件商品的键
//	key2 := "goods" + "541" + ":" + "爱的嫁衣"  //TODO Redis里的第二件商品的键
//	stock, _ := utils.SearchGoodsStock()
//	rdbstock1, _ := utils.RdbLLen(key1)
//	rdbstock2, _ := utils.RdbLLen(key2)
//	// 定义一个每半小时执行一次的定时任务
//	if int64(stock[0]) < rdbstock1 && int64(stock[1]) < rdbstock2 {
//		c.Stop()
//		log.Println("定时任务停止")
//		return
//	}
//	c.AddFunc("*/2 * * * * *", func() {
//
//		log.Println("定时任务54154545454545")
//
//		// 检查库存是否与缓存数据不一致
//		if rdbstock1 != int64(stock[0]) && rdbstock2 != int64(stock[1]) {
//			preheat, _ := utils.PreheatTheProduct()
//
//			// 预热商品库存至缓存
//			for _, g := range *preheat {
//				for i := 0; i < g.GoodsStock; i++ {
//					key := "goods" + strconv.Itoa(int(g.GoodsID)) + ":" + g.GoodsDescription
//					err := utils.RdbLPush(key, g.GoodsDescription)
//					if err != nil {
//						cs.JSON(404, gin.H{"error": "redis同步预热失败" + err.Error()})
//					}
//				}
//				fmt.Println(strconv.Itoa(g.GoodsStock), "][]]]]]]]]]]]]]]]")
//			}
//			cs.JSON(http.StatusOK, "同步预热成功")
//		}
//
//	})
//	// 启动定时任务调度器
//	c.Start()
//
//	// 计时器，一小时后停止定时任务
//	timer := time.NewTimer(1 * time.Hour)
//	<-timer.C
//
//	// 停止定时任务调度器
//	c.Stop()
//}
