var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
import { jsx as _jsx } from "react/jsx-runtime";
import { Component } from "react";
var Clock = /** @class */ (function (_super) {
    __extends(Clock, _super);
    function Clock(props) {
        var _this = _super.call(this, props) || this;
        _this.state = { date: new Date() };
        return _this;
    }
    Clock.prototype.componentDidMount = function () {
        var _this = this;
        this.timer = setInterval(function () { return _this.setState({ date: new Date() }); }, 1000);
    };
    Clock.prototype.componentWillUnmount = function () {
        clearInterval(this.timer);
    };
    Clock.prototype.render = function () {
        return (_jsx("div", { className: "clock", children: _jsx("div", { children: this.state.date.toLocaleTimeString() }) }));
    };
    return Clock;
}(Component));
export default Clock;
