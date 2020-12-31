import {TreeModel} from "./model/TreeModel.js";

"./model/TreeModel.js"

export class Tree {
    constructor() {
        this.name = "tree";
        this.model = new TreeModel();
        this.canvas = document.getElementById(this.name);
        this.ctx = this.canvas.getContext("2d");
        this.backColor = "#4080AA";
        this.foreColor = "#FF0000";
        this.fontColor = "#00C0C0";
    }

    drawBackgroup() {
        this.ctx.fillStyle = this.backColor;
        this.ctx.fillRect(0, 0, this.canvas.offsetWidth, this.canvas.offsetHeight);
    }

    drawTitle() {
        this.ctx.font = "20px Arial";
        this.ctx.fillStyle = this.fontColor;
        var w = (this.canvas.offsetWidth-this.name.length*30)/2;
        this.ctx.fillText(this.model.name, w, 20);
    }

    drawNodes() {

    }

    draw() {
        this.drawBackgroup();
        this.drawTitle();


        this.ctx.beginPath();
        this.ctx.strokeStyle = this.foreColor;
        this.ctx.moveTo(0, 0);
        this.ctx.lineTo(150, 75);
        this.ctx.stroke();
    }

}