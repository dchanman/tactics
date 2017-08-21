window.Util = (function () {
    "use strict";
    return {
        zeropad: function (id) {
            var s = id.toString();
            while (s.length < 6) {
                s = "0" + s;
            }
            return s;
        }
    };
}());