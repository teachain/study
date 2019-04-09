node的配置



//节点私钥

```
func (c *Config) NodeKey() *ecdsa.PrivateKey {
	//如果已经配置过了，直接返回
	// Use any specifically configured key.
	if c.P2P.PrivateKey != nil {
		return c.P2P.PrivateKey
	}
	//如果没有指定数据目录，那么生成私钥并直接返回
	// Generate ephemeral key if no datadir is being used.
	if c.DataDir == "" {
		key, err := crypto.GenerateKey()
		if err != nil {
			log.Crit(fmt.Sprintf("Failed to generate ephemeral node key: %v", err))
		}
		return key
	}
    //从文件中加载
    //路径为${datadir}/geth/nodekey
	keyfile := c.ResolvePath(datadirPrivateKey)
	if key, err := crypto.LoadECDSA(keyfile); err == nil {
		//存在了该文件，则直接返回数据
		return key
	}
	//还没有存储在目录文件中，那么就生成私钥
	// No persistent key found, generate and store a new one.
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Crit(fmt.Sprintf("Failed to generate node key: %v", err))
	}
	//拼接实例路径---${datadir}/geth/
	instanceDir := filepath.Join(c.DataDir, c.name())
	//创建该目录
	if err := os.MkdirAll(instanceDir, 0700); err != nil {
		log.Error(fmt.Sprintf("Failed to persist node key: %v", err))
		return key
	}
	//拼接nodekey的文件路径
	keyfile = filepath.Join(instanceDir, datadirPrivateKey)
	//将nodekey保存到文件中
	if err := crypto.SaveECDSA(keyfile, key); err != nil {
		log.Error(fmt.Sprintf("Failed to persist node key: %v", err))
	}
	//返回nodekey
	return key
}
```



```
// Start create a live P2P node and starts running it.
func (n *Node) Start() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// Short circuit if the node's already running
	if n.server != nil {
		return ErrNodeRunning
	}
	if err := n.openDataDir(); err != nil {
		return err
	}

	// Initialize the p2p server. This creates the node key and
	// discovery databases.
	n.serverConfig = n.config.P2P
	//在这里配置了节点的私钥，它取得的方式看上文的NodeKey方法
	n.serverConfig.PrivateKey = n.config.NodeKey()
	n.serverConfig.Name = n.config.NodeName()
	n.serverConfig.Logger = n.log
	if n.serverConfig.StaticNodes == nil {
		n.serverConfig.StaticNodes = n.config.StaticNodes()
	}
	if n.serverConfig.TrustedNodes == nil {
		n.serverConfig.TrustedNodes = n.config.TrustedNodes()
	}
	if n.serverConfig.NodeDatabase == "" {
		n.serverConfig.NodeDatabase = n.config.NodeDB()
	}
	
	//创建p2p.Server对象
	running := &p2p.Server{Config: n.serverConfig}

	n.log.Info("Starting peer-to-peer node", "instance", n.serverConfig.Name)

	// Otherwise copy and specialize the P2P configuration
	services := make(map[reflect.Type]Service)
	//根据注册进来的构造函数，创建service对象
	for _, constructor := range n.serviceFuncs {
		// Create a new context for the particular service
		ctx := &ServiceContext{
			config:         n.config,
			services:       make(map[reflect.Type]Service),
			EventMux:       n.eventmux,
			AccountManager: n.accman,
		}
		for kind, s := range services { // copy needed for threaded access
			ctx.services[kind] = s
		}
		// Construct and save the service
		service, err := constructor(ctx)
		if err != nil {
			return err
		}
		kind := reflect.TypeOf(service)
		if _, exists := services[kind]; exists {
			return &DuplicateServiceError{Kind: kind}
		}
		services[kind] = service

		log.Error("node start","servive name=",kind)
	}
	//收集协议
	// Gather the protocols and start the freshly assembled P2P server
	for _, service := range services {
		running.Protocols = append(running.Protocols, service.Protocols()...)
	}
	//启动p2p服务器对象
	if err := running.Start(); err != nil {
		return convertFileLockError(err)
	}
	// Start each of the services
	started := []reflect.Type{}
	//启动各项服务
	//这里一个技巧是:它记录了已经启动的服务，然后如果说后续的哪个服务启动失败了
	//就把前面启动的服务全部停止，然后返回错误
	for kind, service := range services {
		// Start the next service, stopping all previous upon failure
		if err := service.Start(running); err != nil {
			for _, kind := range started {
				services[kind].Stop()
			}
			running.Stop()

			return err
		}
		// Mark the service started for potential cleanup
		started = append(started, kind)
	}
	//启动各个服务的rpc接口
	//如果启动过程中遇到错误，也是把启动的服务全部停止，然后返回错误
	//这进一步看到了步步为营的策略。
	// Lastly start the configured RPC interfaces
	if err := n.startRPC(services); err != nil {
		for _, service := range services {
			service.Stop()
		}
		running.Stop()
		return err
	}
	//到这里，就表示p2p服务已经启动成功，服务启动成功，rpc服务启动成功
	// Finish initializing the startup
	n.services = services
	n.server = running
	n.stop = make(chan struct{})

	return nil
}
```

