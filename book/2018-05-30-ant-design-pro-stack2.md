---
layout: post
categories: yield
title: Ant-Design-Pro 基础问题[第四天]
date: 2018-05-30 11:09:03 +0800
description: 一些在群里遇到的问题，以及见解还有解决方案。
keywords: react,ant-design-pro,aliyun
---

其中是我在QQ群 `362329545`遇到的 问题，整理成问答录，希望能解决入门的烦恼，个人见解，如果有那里有错误请指出了我会积极改进，方便他人。



## 问答录[3问]



###  问
   * <a href="#form_1" target="_self">如何编辑修改在一个Form内进行操作</a>
   * <a href="#form_2" target="_self">关于Ant(StandardTable)的默认数据格式</a>
   * <a href="#form_3" target="_self">关于Ant方法按钮方法触发的方式，以及如何定义。</a>
   * <a href="#form_4" target="_self">简单描述一下dispatch调用的过程。</a>
   * <a href="#form_5" target="_self">Redux(Router)页面跳转，携带参数的方式</a>




### 答

####   <span id = "form_1"><font>在Form如何进行编辑以及新增数据。</font></span>

  * 开始这个演示，注意注释描述,代码不一定能运行，只是整体做个介绍

    ```js
        const CreateForm = Form.create()(props => {
        const {modalVisible, form, handleAdd, handleEdit, handleModalVisible, formData, formType} = props;  /** 传入的可变变量 **/
        const okHandle = () => {
            form.validateFields((err, fieldsValue) => {
            if (err) return;
            form.resetFields();
            if (formType === 0) handleAdd(fieldsValue); /**判断类型进行展示数据，以及新增*/
            if (formType === 1) handleEdit(fieldsValue);
            });
        };
        return (
            <Modal
            title="微信配置"
            visible={modalVisible}

            onOk={okHandle} //修改和新增主要看这里  

            onCancel={() => handleModalVisible()}
            >
            //初始化数据 formData

            <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="公钥">
                {form.getFieldDecorator('public_key', {
                initialValue: Object.keys(formData || {}).length ? formData.public_key : '',
                rules: [{required: true, message: 'Please input  public_key...'}],
                })(<Input placeholder="请输入"/>)}
            </FormItem>

            /**这里定义的是 select 的数据以及修改的时候数据的加载 可以保证数据正常显示。**/
            /**目前存在一些疑惑，为什么通过formData 直接绑定，在修改的时候展示的是key 不是 value 通过外部传递就可以。**/
            <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="应用类型">
                {form.getFieldDecorator('channel_way', {
                initialValue: Object.keys(formData || {}).length ? formData.channel_way : '',
                })(
                <Select placeholder="不限" style={{maxWidth: 200, width: '100%'}}>
                    {channelWayData.map(item => (
                    <Option key={item.id} value={item.id}>
                        {item.desc}
                    </Option>
                    ))}
                </Select>
                )}
            </FormItem>
            </Modal>
        );
        });

          handleEdit = fields => {
                message.success(JSON.stringify(fields))
                const {paymentConfigId} = this.state;
                const {dispatch} = this.props;
                fields.channel_id = paymentConfigId;
                dispatch({
                type: 'parkingPayManager/modifyChannel',
                payload: JSON.stringify(fields),
                callback: () => {
                    dispatch({
                    type: 'parkingPayManager/queryChannelList',
                    callback: (e) => {
                        message.success("修改成功，列表已经刷新！")
                    },
                    })
                    this.setState({
                    modalVisible: false, // 清除
                    });
                },
                })
            }

            handleAdd = fields => {
                const {dispatch} = this.props;
                /**switchPress true = 0, fasle = 1 方便数据存储*/ 
                fields.channel_status = this.switchPress(fields.channel_status);
                /**触发action触发model层的state的初始化*/
                dispatch({
                type: 'parkingPayManager/insertChannel',
                payload: JSON.stringify(fields),
                callback: () => {
                    dispatch({
                    type: 'parkingPayManager/queryChannelList',
                    callback: (e) => {
                        message.success('successfully ');
                    },
                    })
                        this.setState({
                            modalVisible: false,
                        });
                    },
                });
            };
        /**传递数据*/
        state = {
            paymentTypeData: [
            {
                id: 1,
                desc: "支付宝",
            }
            ,
            {
                id: 0,
                desc: "微信",
            }
            ,
            ],
            channelWayData: [
            {
                id: 0,
                desc: "小程序,公众号",
            }, {
                id: 1,
                desc: "APP",
            },
            ],
            selectedRows: [],
            formType: 0, /** 0 add  1 edit 定义Form的弹出的模式  **/
            formData: {},
        }
        /**
        * @param  {state:Object} ... 里面多个参数 
        *
        */
        const {isLoading,...[省略]} = this.state;

        const parentMethods = {
            handleAdd: this.handleAdd,
            handleModalVisible: this.handleModalVisible,
            optionChildrenDataTemp: optionChildrenData,
            formData: formData,
            formType: formType,

            handleEdit: this.handleEdit,
            channelWayData: channelWayData,
            paymentTypeData: paymentTypeData,
        };


         /**
         * 组件参数配置（render => （这里面定义的组件））
         * @param {array}  parentMethods 是我们传递的参数组[多个]
         *
         * @param {boolean}  modalVisible 模态框 弹出还暗示隐藏的状态
         */
         <CreateForm {...parentMethods} modalVisible={modalVisible}/>

         /** models 里面的方法 */
         /**
         *@param {Object} payload 是我们提交的参数，或者是用来修改数据等。
         *@param {function} callback 执行成功的回调可以将参数 call回去。
         **/
         effects: {
            * insertChannelCfg({payload, callback}, {call, put}) {
                const response = yield call(insertParkManagerChannelCfg, payload);
                console.log("result Data  successfully");
                console.log(response);
                yield put({
                    type: 'updateStateInsertChannelCfg',
                    payload: response.data,
                    /**
                    *  (response.success) ? response.data : [], 也可以这样去写
                    */
                });
                if (response.success) {
                    callback();
                    /*
                    * callback(response);
                    * 可以将参数携带回去，也可以不携带，具体看下需求
                    **/
                    }
                },
         }
        /**
        *@param {reducers} 通过reducers -> updateStateInsertChannelCfg 将数据传递回去
        **/
        reducers:{
            updateStateInsertChannelCfg(state, action) {
            return {
                ...state,
                InsertChannelCfgNum: action.payload,
                /** InsertChannelCfgNum 记得在 state 里面注册 否则会出现undef */
            }
            },
         }

         /***
         * 解答一下 `namespace: 'parkingPayManager'`
         *  我们在dispatch 里面通过type 调用的 是先通过 namespace 再去找到具体的(effects)内的执行器。
         *
         **/

        /** services [用来发送请求，并且将数据携带回去] **/

        /**
        *@param {JSON} params 我们传递的参数，
        *@func  request 是自己封装的一个请求的方法。
        **/
        export async function insertParkManagerChannelCfg(params) {
            return request('/proxy-post/insertChannelCfg', {
                method: 'POST', data: params || {},
            })
        }
    ```

####   <span id = "form_2"><font>Ant StandardTable 数据格式</font></span>

 * 开始演示

    ```js

        /**
        * 我们通过 model获得的数据将如何渲染出去是个问题，StandardTable是一个 ant封装的组件，如果你要用他的组件就要满足他的数据格式
        **/
            const {parkingPayManager: {channelData}} = this.props;
                const newdata = {
                list: channelData,
                pagination: {
                    total: channelData.length,
                    pageSize: 10,
                    current: 1,
                },
            };

            ...

                <StandardTable
                selectedRows={selectedRows}
                loading={isLoading}
                data={newdata}
                columns={columns}
                />
    ```

####   <span id = "form_3"><font>按钮onClick，渲染触发。</font></span>

* 按钮定义

    ```js
        /**
        * ()=>this.ModifyHandleModalVisible(row,true) 用匿名函数的方式定义不会再加载的时候被触发
        *
        * this.ModifyHandleModalVisible(row,true)  这样会直接被触发（渲染周期开始）。
        **/
        <Button type="edit" onClick={()=>this.ModifyHandleModalVisible(row,true)}>
                修改
        </Button>
    ```
####   <span id = "form_4"><font>dispatch type payload 获得数据的流程(Redux) </font></span>

*  第一个问题什么是 Store , Store主要功能


    1. Store需要负责接收Views传来的Action 然后，根据Action.type和Action.payload对Store里的数据进行修改
    最后，

    2. Store还需要通知Views，数据有改变，Views便去获取最新的Store数据，通过setState进行重新渲染组件（`re-render`）

* Store如何接收来自Views的Action?
    1.   每一个Store实例都拥有dispatch方法，Views只需要通过调用该方法，并传入action对象作为形参，Store自然就就可以收到Action，就像这样：
    ```js
        store.dispatch({
    type: 'modal(namespace)/function Name'
});
    ```
* Store在接收到Action之后，需要根据Action.type和Action.payload修改存储数据，那么，这部分逻辑写在哪里，且怎么将这部分逻辑传递给Store知道呢？
    1. 数据修改逻辑写在Reducer（一个纯函数）里，Store实例在创建的时候，就会被传递这样一个reducer作为形参，这样Store就可以通过Reducer的返回值更新内部数据了，先看一个简单的例子(具体的关于reducer我们后面再讲)：
    ```js
        // 一个reducer
        function counterReducer(state = 0, action) {
        switch (action.type) {
            case 'INCREASE':
            return state + 1;
            case 'DECREASE':
            return state - 1;
            default:
            return state;
        }
        }

        // 传递reducer作为形参
        let store = Redux.createStore(counterReducer);
    ```

* Store通过Reducer修改好了内部数据之后，又是如何通知Views需要获取最新的Store数据来更新的呢？
    1. 每一个Store实例都提供一个subscribe方法，Views只需要调用该方法注册一个回调（内含setState操作），之后在每次dispatch(action)时，该回调都会被触发，从而实现重新渲染；对于最新的Store数据，可以通过Store实例提供的另一个方法getState来获取，就像下面这样：

    ```js
    let unsubscribe = store.subscribe(() =>
        console.log(store.getState())
    );
    //如下

    function createStore(reducer, initialState, enhancer) {
        var currentReducer = reducer
        var currentState = initialState
        var listeners = []

        // 省略若干代码
        //...

        // 通过reducer初始化数据
        dispatch({ type: ActionTypes.INIT })

        return {
            dispatch,
            subscribe,
            getState,
            replaceReducer
        }
    }
    ```
* 归纳总结

    1. Store的数据修改，本质上是通过Reducer来完成的。
    2. Store只提供get方法（即getState），不提供set方法，所以数据的修改一定是通过dispatch(action)来完成，即：action -> reducers -> store
    3. Store除了存储数据之外，还有着消息发布/订阅（pub/sub）的功能，也正是因为这个功能，它才能够同时连接着Actions和Views。

        3.1 dispatch方法 对应着 pub

        3.1 subscribe方法 对应着 sub



####   <span id = "form_5"><font>Ant Redux 页面传递参数，多种方式</font></span>


* path 携带参数 `props.params`

    ```js
    /**传递*/
    import { Router,Route,hashHistory} from 'react-router';
    class App extends React.Component {
        render() {
        ...
        /**使用方式*/
        var data = {id:1,name:1,age:1};
        data = JSON.stringify(data);
        var path = `/user/${data}`;
        ...
        return (
            <Router history={hashHistory}>
                <Route path='/user/:name' component={UserPage}></Route>
            </Router>
        )
        }
    }
    /**接收**/

    export default class UserPage extends React.Component{
        constructor(props){
            super(props);
        }
        render(){
            return(<div>this.props.params.name</div>)
        }
    }
    ```
*   query 类似`request`的`get`


    ```js
        <Route path='/user' component={UserPage}></Route>

       ...
       /**使用方式**/
       var data = {id:1,name:1,age:1};
        var path = {
        pathname:'/user',
        query:data,
        }
       ...
        /**获取数据**/
        var data = this.props.location.query;
        var {id,name,age} = data;
    ```
*  state方式类似于post方式，使用方式和query类似，首先定义路由：

    ```js
        <Route path='/user' component={UserPage}></Route>
        ...
        /**使用*/
        var data = {id:3,name:sam,age:36};
        var path = {
            pathname:'/user',
            state:data,
        }
        ...
        /**获取数据*/
        var data = this.props.location.state;
        var {id,name,age} = data;
    ```
* `state`方式依然可以传递任意类型的数据，而且可以不以明文方式传输.

![](http://47.93.5.46/1.jpeg)


转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
