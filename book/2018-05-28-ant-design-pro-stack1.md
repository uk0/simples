---
layout: post
categories: react
title: ant-design-pro 初级问题以及解答
date: 2018-05-28 15:21:44 +0800
description: ant-design-pro 的基本问题解析
keywords: React
---

## 本章介绍的是ant-design-pro 使用中遇到的一些问题 以及个人的解决方式，大牛勿喷。



## React 与 Ant-Design-Pro


* 下载工程：

```bash
https://github.com/ant-design/ant-design-pro
```
* 安装依赖包

```bash
 npm i -dd  # 等待的时间很长。【可以代理出去速度更快】

 npm start # 启动

 npm run build:aot:prod # 编译

```

* 修改配置文件 `.eslintrc.js` WebStorm 最新版2018 可能会报错，所以需要修改`.eslintrc.js`

```js
module.exports = {
  parser: 'babel-eslint',
  extends: ['airbnb', 'prettier'],
  env: {
    browser: true,
    node: true,
    es6: true,
    mocha: true,
    jest: true,
    jasmine: true,
  },
  rules: {
    // 布尔值类型的 propTypes 的 name 必须为 is 或 has 开头
    // @off 不强制要求写 propTypes
    "react/boolean-prop-naming": "off",
    // 一个 defaultProps 必须有对应的 propTypes
    // @off 不强制要求写 propTypes
    "react/default-props-match-prop-types": "off",
    // 组件必须有 displayName 属性
    // @off 不强制要求写 displayName
    "react/display-name": "off",
    // 禁止在自定义组件中使用一些指定的 props
    // @off 没必要限制
    "react/forbid-component-props": "off",
    // 禁止使用一些指定的 elements
    // @off 没必要限制
    "react/forbid-elements": "off",
    // 禁止使用一些指定的 propTypes
    // @off 不强制要求写 propTypes
    "react/forbid-prop-types": "off",
    // 禁止直接使用别的组建的 propTypes
    // @off 不强制要求写 propTypes
    "react/forbid-foreign-prop-types": "off",
    // 禁止使用数组的 index 作为 key
    // @off 太严格了
    "react/no-array-index-key": "off",
    // 禁止使用 children 做 props
    "react/no-children-prop": "error",
    // 禁止使用 dangerouslySetInnerHTML
    // @off 没必要限制
    "react/no-danger": "off",
    // 禁止在使用了 dangerouslySetInnerHTML 的组建内添加 children
    "react/no-danger-with-children": "error",
    // 禁止使用已废弃的 api
    "react/no-deprecated": "error",
    // 禁止在 componentDidMount 里面使用 setState
    // @off 同构应用需要在 didMount 里写 setState
    "react/no-did-mount-set-state": "off",
    // 禁止在 componentDidUpdate 里面使用 setState
    "react/no-did-update-set-state": "error",
    // 禁止直接修改 this.state
    "react/no-direct-mutation-state": "error",
    // 禁止使用 findDOMNode
    "react/no-find-dom-node": "error",
    // 禁止使用 isMounted
    "react/no-is-mounted": "error",
    // 禁止在一个文件创建两个组件
    // @off 有一个 bug https://github.com/yannickcr/eslint-plugin-react/issues/1181
    "react/no-multi-comp": "off",
    // 禁止在 PureComponent 中使用 shouldComponentUpdate
    "react/no-redundant-should-component-update": "error",
    // 禁止使用 ReactDOM.render 的返回值
    "react/no-render-return-value": "error",
    // 禁止使用 setState
    // @off setState 很常用
    "react/no-set-state": "off",
    // 禁止拼写错误
    "react/no-typos": "error",
    // 禁止使用字符串 ref
    "react/no-string-refs": "error",
    // 禁止在组件的内部存在未转义的 >, ", " 或 }
    "react/no-unescaped-entities": "error",
    // @fixable 禁止出现 HTML 中的属性，如 class
    "react/no-unknown-property": "error",
    // 禁止出现未使用的 propTypes
    // @off 不强制要求写 propTypes
    "react/no-unused-prop-types": "off",
    // 定义过的 state 必须使用
    // @off 没有官方文档，并且存在很多 bug： https://github.com/yannickcr/eslint-plugin-react/search?q=no-unused-state&type=Issues&utf8=%E2%9C%93
    "react/no-unused-state": "off",
    // 禁止在 componentWillUpdate 中使用 setState
    "react/no-will-update-set-state": "error",
    // 必须使用 Class 的形式创建组件
    "react/prefer-es6-class": [
      "error",
      "always"
    ],
    // @fixable props 与 value 之间的等号前后禁止有空格
    "react/jsx-equals-spacing": [
      "error",
      "never"
    ],
    // 数组中的 jsx 必须有 key
    "react/jsx-key": "error",
    // @fixable 限制每行的 props 数量
    // @off 没必要限制
    "react/jsx-max-props-per-line": "off",
    // jsx 中禁止使用 bind
    // @off 太严格了
    "react/jsx-no-bind": "off",
    // 禁止在 jsx 中使用像注释的字符串
    "react/jsx-no-comment-textnodes": "error",
    // 禁止在 jsx 中出现字符串
    // @off 没必要限制
    "react/jsx-no-literals": "off",
    // 禁止使用 target="_blank"
    // @off 没必要限制
    "react/jsx-no-target-blank": "off",
    // 禁止使用未定义的 jsx elemet
    "react/jsx-no-undef": "error",
    // 禁止使用 pascal 写法的 jsx，比如 <TEST_COMPONENT>
    "react/jsx-pascal-case": "error",
    // @fixable props 必须排好序
    // @off 没必要限制
    "react/jsx-sort-props": "off",
    // 限制文件后缀
    // @off 没必要限制
    "react/jsx-filename-extension": "off",

    "import/prefer-default-export": [
      0
    ],
    // 禁止出现重复的 props
    "react/jsx-no-duplicate-props": "error",
    // @fixable jsx 的开始和闭合处禁止有空格
    "prefer-destructuring":[0],
    "react/jsx-tag-spacing":[0],
    "object-shorthand":[0],
    // "react/jsx-tag-spacing": [
    //   "error",
    //   {
    //     "closingSlash": "never",
    //     "beforeSelfClosing": "always",
    //     "afterOpening": "never"
    //   }
    // ],
    // jsx 文件必须 import React
    "react/jsx-uses-react": "error",
    // 定义了的 jsx element 必须使用
    "react/jsx-uses-vars": "error",
    // @fixable 多行的 jsx 必须有括号包起来
    // @off 没必要限制
    "react/jsx-wrap-multilines": "off",
    // @fixable 第一个 prop 必须得换行
    // @off 没必要限制
    "react/jsx-first-prop-new-line": "off",
    // handler 的名称必须是 onXXX 或 handleXXX
    // @off 没必要限制
    "react/jsx-handler-names": "off",
    // 必须使用 pure function
    // @off 没必要限制
    "react/prefer-stateless-function": "off",
    // 组件必须写 propTypes
    // @off 不强制要求写 propTypes
    "react/prop-types": "off",
    // 出现 jsx 的地方必须 import React
    // @off 已经在 no-undef 中限制了
    "react/react-in-jsx-scope": "off",
    // 非 required 的 prop 必须有 defaultProps
    // @off 不强制要求写 propTypes
    "react/require-default-props": "off",
    // 组件必须有 shouldComponentUpdate
    // @off 没必要限制
    "react/require-optimization": "off",
    // render 方法中必须有返回值
    "react/require-render-return": "error",
    // @fixable 组件内没有 children 时，必须使用自闭和写法
    // @off 没必要限制
    "react/self-closing-comp": "off",
    // @fixable 组件内方法必须按照一定规则排序
    "react/sort-comp": "error",
    // propTypes 的熟悉必须按照字母排序
    // @off 没必要限制
    "react/sort-prop-types": "off",
    // style 属性的取值必须是 object
    "react/style-prop-object": "error",
    // HTML 中的自闭和标签禁止有 children
    "react/void-dom-elements-no-children": "error",
    // @fixable 布尔值的属性必须显式的写 someprop={true}
    // @off 没必要限制
    "react/jsx-boolean-value": "off",
    // @fixable 自闭和标签的反尖括号必须与尖括号的那一行对齐
    "react/jsx-closing-bracket-location": [
      "error",
      {
        "nonEmpty": false,
        "selfClosing": "line-aligned"
      }
    ],
    // @fixable 结束标签必须与开始标签的那一行对齐
    // @off 已经在 jsx-indent 中限制了
    "react/jsx-closing-tag-location": "off",
    'generator-star-spacing': [0],
    'consistent-return': [0],
    'no-useless-concat': "off",
    // 'react/forbid-prop-types': [0],
    // 'react/jsx-filename-extension': [1, { extensions: ['.js'] }],
    'global-require': [1],
    // 'import/prefer-default-export': [0],
    // 'react/jsx-no-bind': [0],
    // 'react/prop-types': [0],
    // 'react/prefer-stateless-function': [0],
    // 'react/jsx-wrap-multilines': [
    //   'error',
    //   {
    //     declaration: 'parens-new-line',
    //     assignment: 'parens-new-line',
    //     return: 'parens-new-line',
    //     arrow: 'parens-new-line',
    //     condition: 'parens-new-line',
    //     logical: 'parens-new-line',
    //     prop: 'ignore',
    //   },
    // ],
    'no-else-return': [0],
    'no-restricted-syntax': [0],
    'import/no-extraneous-dependencies': [0],
    'no-use-before-define': [0],
    'jsx-a11y/no-static-element-interactions': [0],
    'jsx-a11y/no-noninteractive-element-interactions': [0],
    'jsx-a11y/click-events-have-key-events': [0],
    'jsx-a11y/anchor-is-valid': [0],
    'no-nested-ternary': [0],
    'prefer-template': [0],
    'no-unused-vars': [0],
    'no-param-reassign': [0],
    'react/sort-comp': [0],
    'import/no-duplicates': [0],
    'no-unneeded-ternary': [0],
    'arrow-body-style': [0],
    'spaced-comment': [0],
    'import/extensions': [0],
    'no-bitwise': [0],
    'no-cond-assign': [0],
    'import/no-unresolved': "off",
    'import/no-useless-concat': [0],
    'comma-dangle': [
      'error',
      {
        arrays: 'always-multiline',
        objects: 'always-multiline',
        imports: 'always-multiline',
        exports: 'always-multiline',
        functions: 'ignore',
      },
    ],
    'object-curly-newline': [0],
    'function-paren-newline': [0],
    'no-restricted-globals': [0],
    'require-yield': [1],
  },
  parserOptions: {
    ecmaFeatures: {
      experimentalObjectRestSpread: true,
    },
  },
  settings: {
    polyfills: ['fetch', 'promises'],
  },
};
```

* 修改默认的 `request.js > request2.js`

* `npm i axios --save -dd `

```javascript

import axios from 'axios';
import { notification } from 'antd';

function checkStatus(response) {
  const { data: result } = response;
  if (response != null) {
      result.success = true;  // get response successfully
      return result;
    }else {
      const error = new Error(result);
      result.success = false;
      error.result = result;
      throw error;
  }
}

/**
 * Requests a URL, returning a promise.
 *
 * @param  {string} url       The URL we want to request
 * @param  {object} [options] The options we want to pass to "axios"
 * @return {object}           An object containing either "data" or "err"
 */
export default function request(url, options) {
  const defaultOptions = {
    credentials: 'include',
  };
  const newOptions = { ...defaultOptions, ...options };

  if (newOptions.method === 'POST' || newOptions.method === 'PUT' || newOptions.method === 'DELETE') {
    if (!(newOptions.body instanceof FormData)) {
      newOptions.headers = {
        Accept: 'application/json',
        'Content-Type': 'application/json; charset=utf-8',
        ...newOptions.headers,
      };
      // newOptions.body = JSON.stringify(newOptions.body);
    } else {
      console.log("----------------request body json-------------------")
      // newOptions.body is FormData
      newOptions.headers = {
        "content-type": "application/json",
        ...newOptions.headers,
      };
    }
  }

  return axios.create().request({
    url,
    method: options && options.method ? options.method : 'GET',
    timeout: 15000, // http请求超时时间
    ...newOptions,
  })
    .then(checkStatus)
    .catch((error) => {
      if (error.code) {
        notification.error({
          message: error.name,
          description: error.message,
        });
      }
      // http请求超时处理
      if ('stack' in error && 'message' in error) {
        const { message } = error;
        if (~message.indexOf('timeout')) {
          notification.error({
            message: `请求错误: ${url}`,
            description: '很抱歉您的请求已经超时了，请稍后再试！',
          });
        } else {
          notification.error({
            message: `请求错误: ${url}`,
            description: error.message,
          });
        }
      }
      const result = { success: false };
      return result;
    });
}
```
* `.roadhogrc.mock.js` 修改如下内容。

```js
import mockjs from 'mockjs';
import { getRule, postRule } from './mock/rule';
import { getActivities, getNotice, getFakeList } from './mock/api';
import { getFakeChartData } from './mock/chart';
import { getProfileBasicData } from './mock/profile';
import { getProfileAdvancedData } from './mock/profile';
import { getNotices } from './mock/notices';
import { format, delay } from 'roadhog-api-doc';

// 代码中会兼容本地 service mock 以及部署站点的静态数据
const proxy = {
  // 支持值为 Object 和 Array
  'GET /proxy-get/(.*)': 'http://127.0.0.1:8004/pay-web-rest/',
  'POST /proxy-post/(.*)': 'http://127.0.0.1:8004/pay-web-rest/',
  'DELETE /proxy-delete/(.*)': 'http://127.0.0.1:8004/pay-web-rest/',
};

// 是否禁用代理
const noProxy = process.env.NO_PROXY === 'true';
export default noProxy ? {} : delay(proxy, 1000);
```

## 写一个模块。

  * 经过上面的配置，现在已经可以开始写模块了。
  * 写一个模块需要配置如下几项。

    ```bash
    ant-design-pro/src/common/menu.js # 页面菜单，可以配置权限。
    ant-design-pro/src/common/router.js # 该文件内配置全局的 route路径，可以分配访问[models]的权限
    ant-design-pro/src/routes # 每一个单独的 route 功能的实现（同等于Page）
    ant-design-pro/src/models # effects 用来注册请求，访问请求。
    ant-design-pro/src/services # 封装 request请求 供给 ‘effects’
    
    ```

   * 现在需要在 `route` 内加一个路径，也就是我们自己的route
   ```js
    // park 服务所支持的 route------start
    '/payList/agent-list': {
      component: dynamicWrapper(app, ['parkingPayManager'], () => import('../routes/PayManager/PayAgent')),
    },
    // 1. 其中 parkingPayManager 是一个models我们要将这个models注册进来，给这个route使用。
    // 2. /payList/agent-list 是`menu`的菜单路径。
   ```
   * 现在需要在 `menu` 内加一个菜单。

   ```js
{
    name: '支付商户管理',
    icon: 'table',
    path: 'payList',
    children: [
      {
        name: '商户管理',
        path: 'agent-list',
      }
    ],
  },
   ```

  * 添加一个model

  ```js
import {
  queryParkManagerAgents,
} from '../services/parkingServices';

export default {
  namespace: 'parkingPayManager', //我们所有的 dispatch 都是通过这个namespace定位过来的。 所以当请求不通，要注意。
  state: {
    QueryAgentData: [],
  },
  effects: {
    *queryAgentList({callback}, { call, put }) { //  `{callback}` 替换成 `_` 就是没有参数
      const response = yield call(queryParkManagerAgents); //queryParkManagerAgents 是定义的service 里面放出的可调用的方法
      console.log("result Data  successfully");
      console.log(response);
      yield put({
        type: 'updateStateQueryAgentList',  // (response.success) ? response.data : [],
        payload: response.data,
      });
      if(response.success)callback(response.data); //进行回调
    },

    //updateStateQueryAgentList 是通过type 进行调用. 名字要相互对应，后面会讲为什么这么调用。
  reducers: {
    updateStateQueryAgentList(state, action) {
      console.log(action.payload);
      return {
        ...state,
        QueryAgentData: action.payload, // 这个QueryAgentData 需要提前定义在state里面
      }
    },
  },
}

  ```

  * 定义一个 `Service`

  ```js

import request from '../utils/request2';

export async function queryParkManagerAgents() {
  return request('/proxy-get/queryAllAgent', { //其中这里面的 proxy-get 是通过 调用mock 的proxy 完成的。 也就是第一次配置的哪个 mock配置文件。
    method: 'GET',
  })
}

  ```

 * 最后编写一个 route （page）

 ```bash
import React, {PureComponent} from 'react';
import {connect} from 'dva';
import {
  Card,
  Button,
  Form,
  Input,
  Select,
  Modal,
  message,
} from 'antd';
import StandardTable from 'components/StandardTable';
import PageHeaderLayout from '../../layouts/PageHeaderLayout';

import styles from './PayAgent.less';

const {Option} = Select;

const FormItem = Form.Item;
//定义一个 Form
const CreateForm = Form.create()(props => {
  const {modalVisible, form, handleAdd, handleModalVisible, optionChildrenDataTemp} = props;
  const okHandle = () => {
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      form.resetFields();
      handleAdd(fieldsValue);
    });
  };
  return (
    <Modal
      title="新建配置"
      visible={modalVisible}
      onOk={okHandle}
      onCancel={() => handleModalVisible()}
    >
      <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="商户代码">
        {form.getFieldDecorator('mch_code', {
          rules: [{required: true, message: 'Please input  mch_code...'}],
        })(<Input placeholder="目前有用"/>)}
      </FormItem>
      <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="通道支付码">
        {form.getFieldDecorator('mchp_code', {
          rules: [{required: true, message: 'Please input  mchp_code...'}],
        })(<Input placeholder="请输入：通道支付码"/>)}
      </FormItem>
      <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="商户名称">
        {form.getFieldDecorator('business_name', {
          rules: [{required: true, message: 'Please input  business_name...'}],
        })(<Input placeholder="请输入： 商户名称"/>)}
      </FormItem>
      <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="描述信息">
        {form.getFieldDecorator('mch_description', {
          rules: [{required: true, message: 'Please input  certificate_file...'}],
        })(<Input placeholder="请输入：描述信息"/>)}
      </FormItem>
      <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="商户ID">
        {form.getFieldDecorator('mch_id', {
          rules: [{required: true, message: 'Please input  mch_id...'}],
        })(<Input placeholder="请输入：商户ID"/>)}
      </FormItem>
      <FormItem labelCol={{span: 5}} wrapperCol={{span: 15}} label="选择通道">
        {form.getFieldDecorator('mchc_code', {
          rules: [{required: true, message: 'Please selectd...'}],
        })(<Select
          // mode="tags"  // multiple  tags combobox
          style={{maxWidth: 200, width: '100%'}}
          placeholder="不限"
        >
          {optionChildrenDataTemp.map(item => (
            <Option key={item.channel_id} value={item.channel_id}>
              {item.channel_name}
            </Option>
          ))}
        </Select>)}
      </FormItem>
    </Modal>
  );
});
// 连接model层的state数据，然后通过this.props.state : loading.models.名(namespace)访问model层的state数据
@connect(({parkingPayManager, loading}) => ({
  parkingPayManager, //获得整个 Model
  loading: loading.models.parkingPayManager, //   true or false 
}))
export default class PayAgent extends PureComponent {
  state = {
    selectedRows: [],
    isLoading: false,
    modalVisible: false,
    optionChildrenData: [],
  };

  componentWillMount() {
    const {dispatch} = this.props;
    dispatch({
      type: 'parkingPayManager/queryChannelList', //通过dispatch调用这个请求。
      callback: (e) => { //请求内的callback  可以做刷新动作等，
        message.success("select init success " + e);
        this.setState({
          optionChildrenData: e,
        });
      },
    });
    console.log("初始化 => 下拉列表的数据");
  }

  componentDidMount() {
    // componentWillMount -> render ->  componentDidMount  ..-> componentWillReceiveProps
    const {dispatch} = this.props;
    // 触发action触发model层的state的初始化
    dispatch({
      type: 'parkingPayManager/queryAgentList',
    });
  }

  handleModalVisible = flag => {
    this.setState({
      modalVisible: !!flag,
    });
  }

  handleAdd = fields => {
    const {dispatch} = this.props;
    // 触发action触发model层的state的初始化
    dispatch({
      type: 'parkingPayManager/insertAgent',
      payload: JSON.stringify(fields),
      callback: () => {
        dispatch({
          type: 'parkingPayManager/queryAgentList',
        });
        this.setState({
          modalVisible: false, // 清除
        });
        message.success('successfully ' + JSON.stringify(fields));
      },
    });
  };

  render() {
    const {parkingPayManager: {QueryAgentData}} = this.props;
    const newdata = {
      list: QueryAgentData,
      pagination: {
        total: QueryAgentData.length,
        pageSize: 10,// 页面量
        current: 1, // 当前第几页
      },
    };
    const {isLoading, modalVisible, selectedRows, optionChildrenData} = this.state;
    const columns = [
      {
        title: '商户伪码',
        dataIndex: 'mch_code',
        align: 'right',
        render: text =>
          text ? (
            text
          ) : (
            '无法显示'
          ),
      },
      {
        title: '商户名称',
        dataIndex: 'business_name',
      },
      {
        title: '描述信息',
        dataIndex: 'mch_description',
      },
      {
        title: '商户的通道码',
        dataIndex: 'mchc_code',
      },
      {
        title: '商户支付码',
        dataIndex: 'mchp_code',
      },
      {
        title: '商户ID',
        dataIndex: 'mch_id',
        render: text =>
          text ? (
            text
          ) : (
            'AliPay Not MchID'
          ),
      },
    ];

    const parentMethods = {
      handleAdd: this.handleAdd,
      handleModalVisible: this.handleModalVisible,
      optionChildrenDataTemp: optionChildrenData,
    };

    return (
      <PageHeaderLayout title="商户配置信息">
        <Card bordered={false}>
          <div className={styles.tableList}>
            <div className={styles.tableListOperator}>
              <Button icon="plus" type="primary" onClick={() => this.handleModalVisible(true)}>
                添加商户配置
              </Button>
            </div>
            <StandardTable
              selectedRows={selectedRows}
              loading={isLoading}
              data={newdata}
              columns={columns}
            />
          </div>
        </Card>
        <CreateForm {...parentMethods} modalVisible={modalVisible}/>
      </PageHeaderLayout>
    );
  }
}

 ```

* 最终完成这个页面，进行访问查看。

![](http://112firshme11224.test.upcdn.net/demos/104d4e32-f0f8-4a54-8f87-a839fce15176.png)





转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
