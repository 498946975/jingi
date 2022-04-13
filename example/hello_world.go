package main

import (
	"fmt"
	"log"
	"runtime"

	"tekao.net/jnigi"
)

func main() {
	//j_path := "/Library/Java/JavaVirtualMachines/jdk1.8.0_271.jdk/Contents/Home/jre/lib/server/libjvm.dylib"
	fmt.Println(jnigi.AttemptToFindJVMLibPath())
	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		log.Fatal(err)
	}
	runtime.LockOSThread()
	// 可以设置-Djava.class.path到jar包解压后的文件夹；
	// 或者设置-Xbootclasspath/a:...为jar包路径，两者应该是等效的
	jvm, env, err := jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, jnigi.DEFAULT_VERSION,
		//[]string{"-Djava.class.path=/Users/oliver/Documents/Lenovo/jar/mqtt-daemon","-verbose:class","-XX:+TraceClassLoading"}))
		//[]string{"-Djava.class.path=/Users/oliver/Documents/Lenovo/jar/mqtt-daemon"}))
		//[]string{"-Djava.class.path=/Users/oliver/Documents/Lenovo/common-utils"}))  // pm
		//[]string{"-Xcheck:jni"}))
		//[]string{"-Xbootclasspath/a:/Users/oliver/Documents/Lenovo/jar/mqtt-daemon.jar"}))
		[]string{"-Xbootclasspath/a:/Users/oliver/Documents/Lenovo/mqtt-daemon.jar:" +
			"/Users/oliver/Documents/Lenovo/mqtt-daemon-extender.jar:" +
			"/Users/oliver/Documents/Lenovo/common-utils-1.0.0.jar:" +
			"/Users/oliver/Documents/Lenovo/common-utils.jar"}))
	if err != nil {
		log.Fatal(err)
	}
	//var aaa jnigi.ObjectRef

	// 1. new pm obj；后续作为参数
	// 1.1 new string, 作为pm的参数
	config, err := env.NewObject("java/lang/String", []byte("/Users/oliver/Desktop/mqtt-daemon.properties"))
	pm, err := env.NewObject("com/lenovo/lrec/cloud/utils/properties/PropertiesManager", config)

	// 2. new string参数；后续作为参数；
	topicId, err := env.NewObject("java/lang/String", []byte("sys/#"))
	// 3. new enginecallbak obj, 1/2作为参数传递
	enginecallbak, err := env.NewObject("com/lenovo/sliic/cloud/mqtt/aed/callback/AedSNMessageReceiveCallBack", topicId, pm)
	fmt.Println(enginecallbak)
	//oarg, err := env.NewObject("java/lang/String", []byte("/Users/oliver/Documeoarg1, err := env.NewObject("java/lang/String", []byte("sys/#"))
	//test, err := env.NewObject("com/lenovo/sliic/cloud/utils/properties/PropertiesManager", oarg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// 当前需要解决的问题：找不到依赖的class文件：com.lenovo.sliic.cloud.mqtt.callback.v1.MessageReceiveToQueueCallBack
	//test, err := env.NewObject("com/lenovo/lrec/cloud/utils/mqtt/callback/BaseMessageCallBack") 					// 成功的
	//test, err := env.NewObject("com/lenovo/lrec/cloud/utils/properties/PropertiesManager", oarg) // pm
	//test, err := env.NewObject("com/lenovo/sliic/cloud/mqtt/callback/MqttEngineCallBack") 								// MqttEngineCallBack
	//test, err := env.NewObject("com/lenovo/sliic/cloud/mqtt/aed/callback/AedSNMessageReceiveCallBack",oarg1, oarg) //AED
	// AedSNMessageReceiveCallBack-->MessageReceiveCallBackV2-->MqttEngineCallBack-->BaseMessageCallBack
	if err != nil {
		log.Fatal(err)
	}

	hello, err := env.NewObject("java/lang/String", []byte("Hello "))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hello)
	world, err := env.NewObject("java/lang/String", []byte("World!"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(world)
	greeting := jnigi.NewObjectRef("java/lang/String")
	err = hello.CallMethod(env, "concat", greeting, world)
	if err != nil {
		log.Fatal(err)
	}

	var goGreeting []byte
	err = greeting.CallMethod(env, "getBytes", &goGreeting)
	if err != nil {
		log.Fatal(err)
	}

	// Prints "Hello World!"
	fmt.Printf("%s\n", goGreeting)

	if err := jvm.Destroy(); err != nil {
		log.Fatal(err)
	}
}
