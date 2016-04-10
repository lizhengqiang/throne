(function () {
    window.Stages = {
        Set: 1,
        Move: 2,
        Event: 0
    };
    window.StageFilter = function (value) {
        switch (value) {
            case Stages.Set:
                return "放置";
            case Stages.Move:
                return "移动";
            case Stages.Event:
                return "事件";
        }
    };
    window.Soldiers = {
        Ship: 0,
        Cavalry: 1,
        Foot: 2,
    };

    window.SoldierFilter = function (soldier) {
        switch (soldier.Type) {
            case Soldiers.Ship:
                return "船" + (soldier.Alive ? "" : "(死亡)");
            case Soldiers.Cavalry:
                return "骑兵" + (soldier.Alive ? "" : "(死亡)");
            case Soldiers.Foot:
                return "步兵" + (soldier.Alive ? "" : "(死亡)");
        }
    };
    window.Commands = {
        Attack: 0,
        Defend: 1,
        Help: 2,
        Money: 3,
        Steal: 4,

    };

    window.CommandFilter = function (command) {
        if (command === null) {
            return "无";
        }
        switch (command.Type) {
            case Commands.Attack:
                return "进攻";
            case Commands.Defend:
                return "防守";
            case Commands.Help:
                return "支援";
            case Commands.Money:
                return "巩固";
            case Commands.Steal:
                return "偷袭";
        }
    };

    window.Resources = {
        Town: 0,
        City: 1,
        Supply: 2,
        Money: 3,

    };

    window.ResourceFilter = function (resource) {
        switch (resource.Type) {
            case Resources.Town:
                return "城镇";
            case Resources.City:
                return "要塞";
            case Resources.Supply:
                return "补给";
            case Resources.Money:
                return "巩固权利";
        }
    }

    Vue.filter("command", CommandFilter);
    Vue.filter("stage", StageFilter);
    Vue.filter("resource", ResourceFilter);
    Vue.filter("soldier", SoldierFilter);


})();