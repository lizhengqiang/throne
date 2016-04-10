(function () {
    var throne = new Vue({
        el: '#throne',
        data: {
            Commands: Commands,
            Stages: Stages,
            Soldiers: Soldiers,
            Resources: Resources,

            select: null,

            players: [],
            player: {},

            game: null,
            war: null,

            areas: [],
            Areas: {}
        },
        methods: {
            loadGame: function () {
                $.get("/game", function (data) {
                    throne.game = data;
                    throne.loadPlayers();
                });
            },
            loadPlayers: function () {
                $.get("/players", function (data) {
                    throne.players = data;
                });
            },
            setPlayer: function (name) {
                throne.players.forEach(function (p) {
                    if (p.Name == name) {
                        throne.player = p;
                        throne.loadAreas();
                    }
                })
            },
            inAreas: function (area, areas) {
                for (var index in areas) {
                    if (areas[index].Id === area.Id) {
                        return true;
                    }
                }
                return false;
            },
            loadAreas: function () {
                $.get("/" + throne.player.Name + "/areas", function (data) {
                    var arounds = {};
                    throne.areas = data;
                    for (var aI in data) {
                        var a = data[aI];
                        throne.Areas[a.Id] = a;
                    }
                    map.option.series[0].nodes = [];
                    map.option.series[0].links = [];
                    map.option.series[0].nodes.push({category: 0, name: '-', value: 25, label: '-'})
                    for (var areaIndex in data) {
                        var area = data[areaIndex];
                        var label = area.Name + "\n" + (area.Belong ? area.Belong.Name : "无");
                        label += "\n资源:";
                        for (var areaResourceIndex in area.Resources) {
                            var resource = area.Resources[areaResourceIndex];
                            label += (ResourceFilter(resource) + ";" )
                        }
                        label += "\n部队:";
                        for (var areaSoldierIndex in area.Soldiers) {
                            var soldier = area.Soldiers[areaSoldierIndex];
                            label += (SoldierFilter(soldier) + ";" )
                        }
                        label += "\n指令:";
                        label += (CommandFilter(area.Command) + ";" );
                        map.option.series[0].nodes.push({
                            category: area.Type,
                            name: area.Id,
                            value: 35,
                            label: label
                        })
                        for (var aroundAreaIndex in area.Around) {
                            var aroundArea = area.Around[aroundAreaIndex];
                            if (arounds[area.Id + "-" + aroundArea.Id]) {
                                continue;
                            }
                            arounds[area.Id + "-" + aroundArea.Id] = true;
                            arounds[aroundArea.Id + "-" + area.Id] = true;
                            map.option.series[0].links.push({
                                source: area.Id,
                                target: aroundArea.Id,
                                weight: 1
                            });
                        }
                    }
                    map.chart.setOption(map.option);
                });
            },
            loadWar: function () {
                $.get("/" + throne.player.Name + "/war", function (data) {
                    throne.war = data;
                });
            },
            leave: function (area) {
                $.get("/" + throne.player.Name + "/" + area.Id + "/select", function (data) {
                    throne.select = {
                        area: area,
                        around: data,
                    };
                    throne.loadAreas();
                    throne.loadGame();
                });
                throne.select = area;
            },
            enter: function (area) {
                $.get("/" + throne.player.Name + "/" + throne.select.area.Id + "/attack/" + area.Id, function () {
                    throne.select = null;
                    throne.loadWar();
                    throne.loadAreas();
                    throne.loadGame();
                });
            },

            help: function (target) {
                $.get("/" + throne.player.Name + "/war/help/" + target, function (data) {
                    throne.loadWar();
                    throne.loadGame();
                });
            },
            back: function (area) {
                $.get("/" + throne.player.Name + "/war/back/" + area.Id, function () {
                    throne.loadWar();
                    throne.loadAreas();
                    throne.loadGame();

                });
            },
            put: function (area, cmd) {
                $.get("/" + throne.player.Name + "/set/" + area.Id + "/" + cmd, function () {
                    throne.loadAreas();
                });
            },
            finishPut: function (area) {
                $.get("/" + throne.player.Name + "/set/finish", function () {
                    throne.loadGame();
                });
            }
        }
    });
    window.throne = throne;
})();