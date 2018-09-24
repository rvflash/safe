const debounce = function (fn, time) {
    let timeout;
    return function() {
        let self = this;
        const fnCall = function () {
            fn.apply(self, arguments);
        }
        clearTimeout(timeout);
        timeout = setTimeout(fnCall, time);
    };
}

Vue.component('vaults', {
    delimiters: ['${', '}'],
    data: function () {
        return {
            items: [],
            hits: 0,
            query: "",
        };
    },
    props: {
        tag: {
            type: String,
            required: true,
        },
    },
    watch: {
        query: function () {
            this.search();
        },
        tag: function () {
            this.list();
        }
    },
    created: function () {
        this.search = debounce(this.list, 500);
    },
    filters: {
        pluralize: function (i) {
            if (i === 0) {
                return 'no vault';
            } else if (i === 1) {
                return '1 vault';
            } else {
                return i + ' vaults';
            }
        }
    },
    methods: {
        add: function() {
            this.$emit('vault', this.tag, '', 'create');
        },
        list: function() {
            let self = this
                , req = new XMLHttpRequest()
                , tag = encodeURIComponent(self.tag)
                , prefix = encodeURIComponent(self.query)
            ;
            // Resets
            self.items = [];
            // Retrieves the list of vaults for this tag.
            req.open('GET', `/tags/${tag}/vaults/?prefix=${prefix}`, true);
            req.onload = function() {
                if (this.status === 200) {
                    JSON.parse(req.responseText).forEach(function(item, key) {
                        item.active = key === 0;
                        if (item.active) {
                            self.$emit('vault', self.tag, item.name, 'view');
                        }
                        self.items[key] = item;
                    });
                    self.hits = self.items.length;
                }
            };
            req.send(null);
        },
        open: function(e) {
            e.preventDefault();
            this.$emit('vault', this.tag, e.currentTarget.title, 'view');
            this.activate(e.currentTarget.title);
        },
        activate: function (name) {
            let self = this;
            self.items.forEach(function(item, key) {
                item.active = name === item.name;
                self.$set(self.items, key, item);
            });
        }
    },
    template: `<div class="vaults">
    <div class="d-flex flex-shrink-0">
        <form class="form-inline">
            <div class="form-group">
                <input class="form-control form-control-sm m-sm-2" type="search" placeholder="Search ..." aria-describedby="hits" v-model="query" name="prefix">
                <small id="hits" class="text-muted" v-show="hits > 0">\${ hits | pluralize }</small>
            </div>
        </form>
        <div class="ml-auto">
            <button type="button" class="btn btn-primary btn-sm m-sm-2" v-on:click="add">Add</button>
        </div>
    </div>
    <div class="list-group list-group-flush flex-md-grow-1" style="overflow: auto" v-if="items.length">
        <a href="#"
            v-bind:class="['list-group-item list-group-item-action', item.active ? 'active' : '']"
            v-bind:title="encodeURIComponent(item.name)"
            v-on:click="open"
            v-for="item in items">
            <p class="lead m-0">\${ item.name }</p>
            <p class="font-weight-light m-0">\${ item.user }</p>
        </a>
    </div>
    <p class="alert alert-light" role="alert" v-else>No vault for \${ tag }</p>
</div>`
});
