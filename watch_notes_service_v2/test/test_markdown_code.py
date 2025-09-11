from markdownify import markdownify as md

def html_to_markdown(html_content):
    """
    将HTML转换为Markdown，正确处理代码块
    """
    # 使用code_language参数确保代码块正确渲染
    markdown = md(
        html_content,
        code_language='',  # 空字符串表示使用默认语言标记
        escape_underscores=False,  # 禁用下划线转义
        # escape_asterisks=False,    # 禁用星号转义
    )
    return markdown

# 测试
html = """## 配置环境变量<br>
<br>
```<i><br>
</i>export MVN_HOME=/opt/soft/apache-maven-3.6.0<br>
export PATH=$MVN_HOME/bin:$PATH<br>
export XGB_HOME=/home/download<br>
export PATH=$XGB_HOME:$PATH<br>
export HDFS_LIB_PATH=${XGB_HOME}/xgboost-packages/libhdfs<br>
export LD_LIBRARY_PATH=${XGB_HOME}/xgboost-packages/lib64:$JAVA_HOME/jre/lib/amd64/server:/${XGB_HOME}/xgboost-packages/libhdfs:$LD_LIBRARY_PATH<br>
export HADOOP_HOME=/opt/cloudera/parcels/CDH/lib/hadoop<br>
export HADOOP_COMMON_HOME=$HADOOP_HOME<br>
export HADOOP_HDFS_HOME=/opt/cloudera/parcels/CDH/lib/hadoop-hdfs<br>
export HADOOP_MAPRED_HOME=/opt/cloudera/parcels/CDH/lib/hadoop-yarn<br>
export HADOOP_YARN_HOME=$HADOOP_MAPRED_HOME<br>
export HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop<br>
```<br>"""

result = html_to_markdown(html)
print(result)