function registerMap(throne) {
    var tooltipFunction = function (params, ticket, callback) {
        var id = params["name"];
        var area = throne.Areas[id];
        return "名称:" + area.Name + "<br>" + "所属:" + area.Belong.Name;
    };
    var option = {
        title: {
            text: '世界地图',
            subtext: '',
            x: 'right',
            y: 'bottom'
        },
        tooltip: {
            trigger: 'item',
            formatter: tooltipFunction,
        },
        toolbox: {
            show: true,
            feature: {
                restore: {show: true},
                magicType: {show: true, type: ['force', 'chord']},
                saveAsImage: {show: true}
            }
        },
        legend: {
            x: 'left',
            data: ['陆地', '海洋']
        },
        series: [
            {
                type: 'force',
                name: "地区",
                ribbonType: false,
                categories: [
                    {
                        name: '陆地'
                    },
                    {
                        name: '海洋'
                    }
                ],
                itemStyle: {
                    normal: {
                        label: {
                            show: true,
                            textStyle: {
                                color: '#333'
                            }
                        },
                        nodeStyle: {
                            brushType: 'both',
                            borderColor: 'rgba(255,215,0,0.4)',
                            borderWidth: 1,
                            r: 50
                        },
                        linkStyle: {
                            type: 'curve',
                            width: '2',
                            color: '#000000',
                        }
                    },
                    emphasis: {
                        label: {
                            show: false
                            // textStyle: null      // 默认使用全局文本样式，详见TEXTSTYLE
                        },
                        nodeStyle: {
                            r: 50
                        },
                        linkStyle: {
                            type: 'curve',
                            width: '2',
                            color: '#000000',
                        }
                    }
                },
                useWorker: false,
                minRadius: 25,
                maxRadius: 35,
                gravity: 1.1,
                scaling: 1.1,
                roam: 'move',
                nodes: [
//                        {category: 0, name: '乔布斯', value: 10, label: '乔布斯\n（主要）'},
//                        {category: 1, name: '丽萨-乔布斯', value: 2},
//                        {category: 1, name: '保罗-乔布斯', value: 3},
//                        {category: 1, name: '克拉拉-乔布斯', value: 3},
//                        {category: 1, name: '劳伦-鲍威尔', value: 7},
//                        {category: 2, name: '史蒂夫-沃兹尼艾克', value: 5},
//                        {category: 2, name: '奥巴马', value: 8},
//                        {category: 2, name: '比尔-盖茨', value: 9},
//                        {category: 2, name: '乔纳森-艾夫', value: 4},
//                        {category: 2, name: '蒂姆-库克', value: 4},
//                        {category: 2, name: '龙-韦恩', value: 1},
                ],
                links: [
//                        {source: '丽萨-乔布斯', target: '乔布斯', weight: 1, name: '女儿'},
//                        {source: '保罗-乔布斯', target: '乔布斯', weight: 2, name: '父亲'},
//                        {source: '克拉拉-乔布斯', target: '乔布斯', weight: 1, name: '母亲'},
//                        {source: '劳伦-鲍威尔', target: '乔布斯', weight: 2},
//                        {source: '史蒂夫-沃兹尼艾克', target: '乔布斯', weight: 3, name: '合伙人'},
//                        {source: '奥巴马', target: '乔布斯', weight: 1},
//                        {source: '比尔-盖茨', target: '乔布斯', weight: 6, name: '竞争对手'},
//                        {source: '乔纳森-艾夫', target: '乔布斯', weight: 1, name: '爱将'},
//                        {source: '蒂姆-库克', target: '乔布斯', weight: 1},
//                        {source: '龙-韦恩', target: '乔布斯', weight: 1},
//                        {source: '克拉拉-乔布斯', target: '保罗-乔布斯', weight: 1},
//                        {source: '奥巴马', target: '保罗-乔布斯', weight: 1},
//                        {source: '奥巴马', target: '克拉拉-乔布斯', weight: 1},
//                        {source: '奥巴马', target: '劳伦-鲍威尔', weight: 1},
//                        {source: '奥巴马', target: '史蒂夫-沃兹尼艾克', weight: 1},
//                        {source: '比尔-盖茨', target: '奥巴马', weight: 6},
//                        {source: '比尔-盖茨', target: '克拉拉-乔布斯', weight: 1},
//                        {source: '蒂姆-库克', target: '奥巴马', weight: 1}
                ]
            }
        ]
    };
//        function focus(param) {
//            var data = param.data;
//            var links = option.series[0].links;
//            var nodes = option.series[0].nodes;
//            if (
//                    data.source !== undefined
//                    && data.target !== undefined
//            ) { //点击的是边
//                var sourceNode = nodes.filter(function (n) {
//                    return n.name == data.source
//                })[0];
//                var targetNode = nodes.filter(function (n) {
//                    return n.name == data.target
//                })[0];
//                console.log("选中了边 " + sourceNode.name + ' -> ' + targetNode.name + ' (' + data.weight + ')');
//            } else { // 点击的是点
//                console.log("选中了" + data.name + '(' + data.value + ')');
//            }
//        }

    var myChart = echarts.init(document.getElementById('map'));

    window.map = {
        chart: myChart,
        option: option,
        tooltipFunction: tooltipFunction
    };
};

