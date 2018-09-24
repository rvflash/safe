const app = new Vue({
    delimiters: ['${', '}'],
    data: {
        tag: "",
        vault: "",
        action: "create"
    },
    el: '#app',
    methods: {
        one: function (tag, vault, action) {
            this.tag = tag;
            this.vault = vault;
            this.action = action;
        },
        all: function (tag) {
            this.tag = tag;
        }
    }
});