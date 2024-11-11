---
layout: post
title: Java进阶之 如何自动生成代码[学习]
categories: Java, FreeMarker
description: Java, FreeMarker
keywords: Java, FreeMarker
---

 什么样的场景和代码适合用自动生成这种方式呢？
 做过Java服务端的朋友一定都知道代码中我们需要编写与数据库表映射的Java实体类（Entity）、
 需要编写与实体对应的DAO类（XxDao.java类中有包含对应实体的增、删、改、查基本操作）。在这些实体类中通常都是一些属性方法以及属性对应的get/set方法、而实体对应的DAO类中也基本会包含有增、删、改、查这些与数据库操作相关的方法。在编写了那么多的实体类和Dao类的过程中 你是否发现了这些代码中有很多地方都相似或者差不多、只是名字不同而已呢？对、那么这个时候其实我们可以定义一个模板、通过模板我们来让代码自动生成去吧。

## 一、前言：为什么要有代码的自动生成？
 ```text
  对于这个问题 最简洁直接的回答就是：代替手动编写代码、提高工作效率。
 ```
## 二、FreeMarker的简单介绍
 在进入正文前，让我们首先简单、快速了解一下FreeMarker。
  （做过Web开发的朋友肯定都是相当熟悉的 我当时 也是在做Web开发的时候第一次接触了FreeMarker）
* 1、概述：FreeMarker是一款模板引擎：即一种基于模板、用来生成输出文本的通用工具。更多的是被用来设计生成HTML页面。
   简单说就是：FreeMarker是使用模板生成文本页面来呈现已经准备好的数据。如下图表述
   
![描述](http://112firshme11224.test.upcdn.net/posts/java/20141227152045529.png)

 FreeMarker官网：http://freemarker.org/
* 2通过一个简单的例子来展示如何使用FreeMarker定义模板、绑定模型数据、生成最终显示的Html页面：
  2.1.新建项目 在项目根目录下新建"template"文件夹，用来存放我们的Template file,
   如我们新建模板文件test.ftl  内容如下：
   ```html
   <html>  
   <head>  
   <title>Welcome!</title>  
   </head>  
   <body>  
   <h1>  
   Welcome ${user}<#if user == "Big Joe">, our beloved leader</#if>!  
   </h1>  
   <p>Our latest product:  
   <a href="${latestProduct.url}">${latestProduct.name}</a>!  
   </body>  
   </html>
   ```
   2.2.项目引入freemarker.jar（下载地址：https://jarfiles.pandaidea.com/freemarker.html），
      在Java类中使用FreeMarker API方法引用模板文件、创建数据模型、合并数据模型与模板文件最终输入，
   代码如下：
   ```java
   import java.io.File;  
   import java.io.IOException;  
   import java.io.OutputStreamWriter;  
   import java.io.Writer;  
   import java.util.HashMap;  
   import java.util.Map;  
     
   import freemarker.template.Configuration;  
   import freemarker.template.DefaultObjectWrapper;  
   import freemarker.template.Template;  
   import freemarker.template.TemplateException;  
     
   public class HtmlGeneratorClient {  
     
       public static void main(String[] args) {  
           try {  
               Configuration cfg = new Configuration();  
               /* 指定模板文件从何处加载的数据源，这里设置成一个文件目录   */
               cfg.setDirectoryForTemplateLoading(new File("./template"));  
               cfg.setObjectWrapper(new DefaultObjectWrapper());  
                 
               /* 获取或创建模板   */
               Template template = cfg.getTemplate("test.ftl");  
                 
               /* 创建数据模型   */
               Map root = new HashMap();  
               root.put("user", "Big Joe");          
               Map latest = new HashMap();  
               root.put("latestProduct", latest);  
               latest.put("url", "products/greenmouse.html");  
               latest.put("name", "green mouse");  
                 
               /* 将模板和数据模型合并 输出到Console  */
               Writer out = new OutputStreamWriter(System.out);  
               template.process(root, out);  
               out.flush();  
                 
           } catch (IOException e) {  
               e.printStackTrace();  
           } catch (TemplateException e) {  
               e.printStackTrace();  
           }  
     
       }  
     
   }  
   ```
   
   2.3.最终生成的HTML的页面代码如下：
   ```html
   <html>  
   <head>  
   <title>Welcome!</title>  
   </head>  
   <body>  
   <h1>  
   Welcome Big Joe, our beloved leader!  
   </h1>  
   <p>Our latest product:  
   <a href="products/greenmouse.html">green mouse</a>!  
   </body>  
   </html>
   ```
## 三、如何使用FreeMerker完成Java类代码的自动生成

* 1、属性类型的枚举类 PropertyType.java
    
    ```java
    /** 
     * 属性类型枚举类 
     * @author github.com/zmatsh 
     * 
     */  
    public enum PropertyType {  
        Byte, Short, Int, Long, Boolean, Float, Double, String, ByteArray, Date  
    } 
    ```
* 2、实体对应的字段属性类 Property.java

```java
/** 
 * 实体对应的属性类 
 * @author github.com/zmatsh 
 * 
 */  
public class Property {  
    // 属性数据类型  
    private String javaType;  
    // 属性名称  
    private String propertyName;  
      
    private PropertyType propertyType;  
      
    public String getJavaType() {  
        return javaType;  
    }  
  
    public void setJavaType(String javaType) {  
        this.javaType = javaType;  
    }  
  
    public String getPropertyName() {  
        return propertyName;  
    }  
  
    public void setPropertyName(String propertyName) {  
        this.propertyName = propertyName;  
    }  
  
    public PropertyType getPropertyType() {  
        return propertyType;  
    }  
  
    public void setPropertyType(PropertyType propertyType) {  
        this.propertyType = propertyType;  
    }  
          
}  
```

*  3、实体模型类 Entity.java

```java

import java.util.List;  
  
/** 
 * 实体类 
 * @author github.com/zmatsh 
 * 
 */  
public class Entity {  
    // 实体所在的包名  
    private String javaPackage;  
    // 实体类名  
    private String className;  
    // 父类名  
    private String superclass;  
    // 属性集合  
    List<Property> properties;  
    // 是否有构造函数  
    private boolean constructors;     
      
    public String getJavaPackage() {  
        return javaPackage;  
    }  
      
    public void setJavaPackage(String javaPackage) {  
        this.javaPackage = javaPackage;  
    }  
      
    public String getClassName() {  
        return className;  
    }  
      
    public void setClassName(String className) {  
        this.className = className;  
    }  
      
    public String getSuperclass() {  
        return superclass;  
    }  
      
    public void setSuperclass(String superclass) {  
        this.superclass = superclass;  
    }  
      
    public List<Property> getProperties() {  
        return properties;  
    }  
      
    public void setProperties(List<Property> properties) {  
        this.properties = properties;  
    }  
      
    public boolean isConstructors() {  
        return constructors;  
    }  
      
    public void setConstructors(boolean constructors) {  
        this.constructors = constructors;  
    }     
  
}  
```

*   4、在项目根目录下新建"template"文件夹，用来存放我们的Template file, 新建实体模板entity.ftl 内容如下：

```text
package ${entity.javaPackage};  
  
/**  
 * This code is generated by FreeMarker  
 * @authorgithub.com/zmatsh  
 *  
 */  
public class ${entity.className}<#if entity.superclass?has_content> extends ${entity.superclass} </#if>  
{  
    /********** attribute ***********/  
<#list entity.properties as property>  
    private ${property.javaType} ${property.propertyName};  
      
</#list>  
    /********** constructors ***********/  
<#if entity.constructors>  
    public ${entity.className}() {  
      
    }  
  
    public ${entity.className}(<#list entity.properties as property>${property.javaType} ${property.propertyName}<#if property_has_next>, </#if></#list>) {  
    <#list entity.properties as property>  
        this.${property.propertyName} = ${property.propertyName};  
    </#list>  
    }  
</#if>  
  
    /********** get/set ***********/  
<#list entity.properties as property>  
    public ${property.javaType} get${property.propertyName?cap_first}() {  
        return ${property.propertyName};  
    }  
  
    public void set${property.propertyName?cap_first}(${property.javaType} ${property.propertyName}) {  
        this.${property.propertyName} = ${property.propertyName};  
    }  
      
</#list>  
}  
```

*   5、自动生成实体类 客户端代码 EntityGeneratorClient.java

```java
import java.io.File;  
import java.io.FileWriter;  
import java.io.IOException;  
import java.io.OutputStreamWriter;  
import java.io.Writer;  
import java.util.ArrayList;  
import java.util.HashMap;  
import java.util.List;  
import java.util.Map;  
  
import freemarker.template.Configuration;  
import freemarker.template.DefaultObjectWrapper;  
import freemarker.template.Template;  
import freemarker.template.TemplateException;  
/** 
 * 自动生成实体类客户端 
 * @authorgithub.com/zmatsh 
 * 
 */  
public class EntityGeneratorClient {  
      
    private static File javaFile = null;  
  
    public static void main(String[] args) {  
        Configuration cfg = new Configuration();      
        try {   
            // 步骤一：指定 模板文件从何处加载的数据源，这里设置一个文件目录  
            cfg.setDirectoryForTemplateLoading(new File("./template"));  
            cfg.setObjectWrapper(new DefaultObjectWrapper());  
              
            // 步骤二：获取 模板文件  
            Template template = cfg.getTemplate("entity.ftl");  
              
            // 步骤三：创建 数据模型  
            Map<String, Object> root = createDataModel();  
              
            // 步骤四：合并 模板 和 数据模型  
            // 创建.java类文件  
            if(javaFile != null){  
                Writer javaWriter = new FileWriter(javaFile);  
                template.process(root, javaWriter);  
                javaWriter.flush();  
                System.out.println("文件生成路径：" + javaFile.getCanonicalPath());  
                  
                javaWriter.close();  
            }  
            // 输出到Console控制台  
            Writer out = new OutputStreamWriter(System.out);  
            template.process(root, out);  
            out.flush();  
            out.close();  
              
        } catch (IOException e) {  
            e.printStackTrace();  
        } catch (TemplateException e) {  
            e.printStackTrace();  
        }  
  
    }  
  
      
    /** 
     * 创建数据模型 
     * @return 
     */  
    private static Map<String, Object> createDataModel() {  
        Map<String, Object> root = new HashMap<String, Object>();  
        Entity user = new Entity();  
        user.setJavaPackage("com.study.entity"); // 创建包名  
        user.setClassName("User");  // 创建类名  
        user.setConstructors(true); // 是否创建构造函数  
        // user.setSuperclass("person");  
          
        List<Property> propertyList = new ArrayList<Property>();  
          
        // 创建实体属性一   
        Property attribute1 = new Property();  
        attribute1.setJavaType("String");  
        attribute1.setPropertyName("name");  
        attribute1.setPropertyType(PropertyType.String);  
          
        // 创建实体属性二  
        Property attribute2 = new Property();  
        attribute2.setJavaType("int");  
        attribute2.setPropertyName("age");  
        attribute2.setPropertyType(PropertyType.Int);  
          
        propertyList.add(attribute1);  
        propertyList.add(attribute2);  
          
        // 将属性集合添加到实体对象中  
        user.setProperties(propertyList);  
          
        // 创建.java类文件  
        File outDirFile = new File("./src-template");  
        if(!outDirFile.exists()){  
            outDirFile.mkdir();  
        }  
          
        javaFile = toJavaFilename(outDirFile, user.getJavaPackage(), user.getClassName());  
          
        root.put("entity", user);  
        return root;  
    }  
      
      
    /** 
     * 创建.java文件所在路径 和 返回.java文件File对象 
     * @param outDirFile 生成文件路径 
     * @param javaPackage java包名 
     * @param javaClassName java类名 
     * @return 
     */  
    private static File toJavaFilename(File outDirFile, String javaPackage, String javaClassName) {  
        String packageSubPath = javaPackage.replace('.', '/');  
        File packagePath = new File(outDirFile, packageSubPath);  
        File file = new File(packagePath, javaClassName + ".java");  
        if(!packagePath.exists()){  
            packagePath.mkdirs();  
        }  
        return file;  
    }  
  
}  

```

# 四、背后的思考
  通过上面两个简单的示例我们了解到所谓的自动生成代码其实就是：
  1、定义java类模板文件 2、定义模板数据  3、引用模板文件（.ftl）与模板数据合并生成Java类。
  上面的示例中 有的朋友可能会问不就是要编写一个实体对象吗？干嘛搞那么麻烦、又建.ftl文件、又写了那么多类、定义模板数据的过程也是那么麻烦、我还不如手动去写、声明几个属性、set/get快捷键一下子就编写好啦。 真的是这样吗？
  从一个辅助工具和软件架构的方面去思考，假设做成一个开发的辅助工具或是插件去完成实体类和对应DAO类的自动生成。假设需要建10个实体类和对应含有增删改查基本操作的DAO类。我在C/S客户端上填写包名、类名、属性字段等信息 然后一键生成，想想那是多么爽、多么痛快的一件事（当然 前提是你的模板类要编写的非常强大、通用），而你也许还在不停的 Ctrl+C、Ctrl+V。

