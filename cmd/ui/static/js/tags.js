Vue.component('tags', {
    delimiters: ['${', '}'],
    data: function () {
        return {
            items: [],
        };
    },
    created: function () {
        this.list("");
    },
    methods: {
        add: function () {
            let tag = window.prompt("Create a new tag named:")
            if (!tag) {
                return;
            }
            let req = new XMLHttpRequest();
            req.open('POST', '/tag', false);
            req.setRequestHeader("Content-Type", "application/json");
            req.send(JSON.stringify({name:tag}));
            if (req.status === 200) {
                // Refreshes the list to add the new one.
                this.list(tag);
            } else {
                window.alert("Failed to create the new tag");
            }
        },
        list: function(cur) {
            let self = this
                , req = new XMLHttpRequest();
            req.open('GET', '/tag', true);
            req.onload = function() {
                self.items = [];
                if (this.status === 200) {
                    let res = JSON.parse(req.responseText);
                    res.forEach(function(tag) {
                        if (cur === "") {
                            // By default, the first is the one activated.
                            cur = tag;
                            self.$emit('tag', tag);
                        }
                        self.items.push({ active: cur === tag, name: tag });
                    });
                }
            };
            req.send(null);
        },
        open: function (e) {
            e.preventDefault();
            this.$emit('tag', e.currentTarget.title);
            this.activate(e.currentTarget.title);
        },
        activate: function (name) {
            this.items.forEach(function(item) {
                item.active = name === item.name;
            });
        }
    },
    template: `<div class="tags">
    <div class="clearfix">
        <h6 class="float-left text-white pt-3">Tags</h6>
        <button type="button" class="btn btn-secondary btn-sm my-sm-2 float-right" v-on:click="add">Add</button>
    </div>
    <ul class="nav flex-column nav-pills flex-md-grow-1" style="overflow: auto" v-if="items.length">
        <li class="nav-item" v-for="item in items">
            <a href="#"
                v-bind:class="['nav-link', item.active ? 'active' : '']"
                v-bind:title="encodeURIComponent(item.name)"
                v-on:click="open">\${ item.name }
            </a>
        </li>
    </ul>
    <p class="alert alert-info" role="alert" v-else>No tag yet</p>
</div>`
});




