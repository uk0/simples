---
layout: post
title: Flink Graph Vertex[来自项目中的经历]
categories: Idea,java
description: Flink学习笔记
keywords: Flink, Graph，StreamGraph
---


发这个贴的原因，是因为 是因为 是因为啥来着，忘了 那就是因为你，因为你 。


## 摘要: Graph构建
   * `Graph`对象是用户的操作入口，主要包含`edge`和`vertex`两部分。边是由点组成，所以边中所有的点就是点的全集，但这个全集包含了重复的点，去重后就是`Vertex`。


## Flink Graph Demo

   * ConnectedComponentsData

```java
package com.todcloud.flink.StreamGraph.utils;

/**
 * Created by zhangjianxin on 2017/7/14.
 */
import org.apache.flink.api.java.DataSet;
import org.apache.flink.api.java.ExecutionEnvironment;
import org.apache.flink.api.java.tuple.Tuple2;

import java.util.LinkedList;
import java.util.List;

/**
 * Provides the default data sets used for the Connected Components example program.
 * The default data sets are used, if no parameters are given to the program.
 *
 */
public class ConnectedComponentsData {

    public static final long[] VERTICES  = new long[] {
            1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16};
    /**至高点*/
    public static DataSet<Long> getDefaultVertexDataSet(ExecutionEnvironment env) {
        List<Long> verticesList = new LinkedList<Long>();
        for (long vertexId : VERTICES) {
            verticesList.add(vertexId);
        }
        return env.fromCollection(verticesList);
    }
    /**边缘*/
    public static final Object[][] EDGES = new Object[][] {
            new Object[]{1L, 2L},
            new Object[]{2L, 3L},
            new Object[]{2L, 4L},
            new Object[]{3L, 5L},
            new Object[]{6L, 7L},
            new Object[]{8L, 9L},
            new Object[]{8L, 10L},
            new Object[]{5L, 11L},
            new Object[]{11L, 12L},
            new Object[]{10L, 13L},
            new Object[]{9L, 14L},
            new Object[]{13L, 14L},
            new Object[]{1L, 15L},
            new Object[]{16L, 1L}
    };
    /**制造数据采集*/
    public static DataSet<Tuple2<Long, Long>> getDefaultEdgeDataSet(ExecutionEnvironment env) {

        List<Tuple2<Long, Long>> edgeList = new LinkedList<Tuple2<Long, Long>>();
        for (Object[] edge : EDGES) {
            edgeList.add(new Tuple2<Long, Long>((Long) edge[0], (Long) edge[1]));
        }
        return env.fromCollection(edgeList);
    }

}
```

  * ConnectedComponents

```java
package com.todcloud.flink.StreamGraph;

/**
 * Created by zhangjianxin on 2017/7/14.
 */
import com.todcloud.flink.StreamGraph.utils.ConnectedComponentsData;
import org.apache.flink.api.common.functions.FlatJoinFunction;
import org.apache.flink.api.common.functions.FlatMapFunction;
import org.apache.flink.api.common.functions.JoinFunction;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.api.java.DataSet;
import org.apache.flink.api.java.ExecutionEnvironment;
import org.apache.flink.api.java.aggregation.Aggregations;
import org.apache.flink.api.java.functions.FunctionAnnotation.ForwardedFields;
import org.apache.flink.api.java.functions.FunctionAnnotation.ForwardedFieldsFirst;
import org.apache.flink.api.java.functions.FunctionAnnotation.ForwardedFieldsSecond;
import org.apache.flink.api.java.operators.DeltaIteration;
import org.apache.flink.api.java.tuple.Tuple1;
import org.apache.flink.api.java.tuple.Tuple2;
import org.apache.flink.api.java.utils.ParameterTool;
import org.apache.flink.util.Collector;

/**
 *使用增量迭代实现连接的组件算法
 *
 * <p>最初，算法为每个顶点分配一个惟一的ID. 在每一步中，顶点选择自己ID的最小值
 *邻居的I Ds，作为它的新ID，并告诉它的邻居它的新ID。算法完成后，同一组件中的所有顶点将具有相同的ID
 *
 * <p>一个顶点的组件 ID没有改变，不需要在下一步传播它的信息 的一步。正因为如此,
 */
@SuppressWarnings("serial")
public class ConnectedComponents {


    public static void main(String... args) throws Exception {

        /** 检查输入参数 */
        final ParameterTool params = ParameterTool.fromArgs(args);

        /** 建立执行环境 */
        ExecutionEnvironment env = ExecutionEnvironment.getExecutionEnvironment();

        final int maxIterations = params.getInt("iterations", 10);

        /** 在web接口中提供可用的参数 */
        env.getConfig().setGlobalJobParameters(params);

        /** 顶点和边读数据 */
        DataSet<Long> vertices = getVertexDataSet(env, params);
        DataSet<Tuple2<Long, Long>> edges = getEdgeDataSet(env, params).flatMap(new UndirectEdge());

        /** 分配初始组件(等于顶点id) */
        DataSet<Tuple2<Long, Long>> verticesWithInitialId = vertices.map(new DuplicateValue<Long>());

        /**打开一个增量迭代 */
        DeltaIteration<Tuple2<Long, Long>, Tuple2<Long, Long>> iteration =
                verticesWithInitialId.iterateDelta(verticesWithInitialId, maxIterations, 0);

        /** 应用步骤逻辑 */
        /** 加入边缘，选择最小的邻居 */
        /** 如果候选的组件较小，则更新 */
        DataSet<Tuple2<Long, Long>> changes = iteration.getWorkset().join(edges).where(0).equalTo(0).with(new NeighborWithComponentIDJoin())
                .groupBy(0).aggregate(Aggregations.MIN, 1)
                .join(iteration.getSolutionSet()).where(0).equalTo(0)
                .with(new ComponentIdFilter());

        /**关闭增量迭代(增量和新工作集是相同的) */
        DataSet<Tuple2<Long, Long>> result = iteration.closeWith(changes, changes);

        /**排放结果 */
        if (params.has("output")) {
            result.writeAsCsv(params.get("output"), "\n", " ");

            env.execute("连接组件的DEMO");
        } else {
            System.out.println("打印输出到stdout. Use --output to specify output path.");
            result.print();
        }
    }


    /**
     * 函数将值转换为两个字段的值
     */
    @ForwardedFields("*->f0")
    public static final class DuplicateValue<T> implements MapFunction<T, Tuple2<T, T>> {

        @Override
        public Tuple2<T, T> map(T vertex) {
            return new Tuple2<T, T>(vertex, vertex);
        }
    }

    /**
     *
     通过向每个输入边缘发射输入边本身和一个反向的版本，以无定向的边缘
     */
    public static final class UndirectEdge implements FlatMapFunction<Tuple2<Long, Long>, Tuple2<Long, Long>> {
        Tuple2<Long, Long> invertedEdge = new Tuple2<Long, Long>();

        @Override
        public void flatMap(Tuple2<Long, Long> edge, Collector<Tuple2<Long, Long>> out) {
            invertedEdge.f0 = edge.f1;
            invertedEdge.f1 = edge.f0;
            out.collect(edge);
            out.collect(invertedEdge);
        }
    }

    /**
     * UDF that joins a (Vertex-ID, Component-ID) pair that represents the current component that
     * a vertex is associated with, with a (Source-Vertex-ID, Target-VertexID) edge. The function
     * produces a (Target-vertex-ID, Component-ID) pair.
     */
    @ForwardedFieldsFirst("f1->f1")
    @ForwardedFieldsSecond("f1->f0")
    public static final class NeighborWithComponentIDJoin implements JoinFunction<Tuple2<Long, Long>, Tuple2<Long, Long>, Tuple2<Long, Long>> {

        @Override
        public Tuple2<Long, Long> join(Tuple2<Long, Long> vertexWithComponent, Tuple2<Long, Long> edge) {
            return new Tuple2<Long, Long>(edge.f1, vertexWithComponent.f1);
        }
    }

    /**
     * Emit the candidate (Vertex-ID, Component-ID) pair if and only if the
     * candidate component ID is less than the vertex's current component ID.
     */
    @ForwardedFieldsFirst("*")
    public static final class ComponentIdFilter implements FlatJoinFunction<Tuple2<Long, Long>, Tuple2<Long, Long>, Tuple2<Long, Long>> {

        @Override
        public void join(Tuple2<Long, Long> candidate, Tuple2<Long, Long> old, Collector<Tuple2<Long, Long>> out) {
            if (candidate.f1 < old.f1) {
                out.collect(candidate);
            }
        }
    }

    /**  UTIL 研究方法  */

    private static DataSet<Long> getVertexDataSet(ExecutionEnvironment env, ParameterTool params) {
        if (params.has("vertices")) {
            return env.readCsvFile(params.get("vertices")).types(Long.class).map(
                    new MapFunction<Tuple1<Long>, Long>() {
                        public Long map(Tuple1<Long> value) {
                            return value.f0;
                        }
                    });
        } else {
            System.out.println("Executing Connected Components example with default vertices data set.");
            System.out.println("Use --vertices to specify file input.");
            return ConnectedComponentsData.getDefaultVertexDataSet(env);
        }
    }

    private static DataSet<Tuple2<Long, Long>> getEdgeDataSet(ExecutionEnvironment env, ParameterTool params) {
        if (params.has("edges")) {
            return env.readCsvFile(params.get("edges")).fieldDelimiter(" ").types(Long.class, Long.class);
        } else {
            System.out.println("Executing Connected Components example with default edges data set.");
            System.out.println("Use --edges to specify file input.");
            return ConnectedComponentsData.getDefaultEdgeDataSet(env);
        }
    }
}
```


