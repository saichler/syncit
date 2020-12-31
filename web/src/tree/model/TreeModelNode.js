export class TreeModelNode {
    constructor() {
        this.name = "node"
        this.children = new Array(0);
    }

    isLeaf() {
        return this.children.length === 0;
    }
}