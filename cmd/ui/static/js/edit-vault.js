function editDefaultData() {
    return {
        digits: 10,
        error: '',
        kind: 0,
        length: 64,
        note: '',
        noUpper: false,
        pass: '',
        repeat: true,
        symbols: 10,
        url: '',
        user: '',
    };
}

Vue.component('edit-vault', {
    delimiters: ['${', '}'],
    data: function () {
        return editDefaultData();
    },
    props: {
        id: {
            type: String,
            default: "edit",
        },
        tag: {
            type: String,
            required: true,
        },
        vault: {
            type: String,
            required: true,
        },
        action: {
            type: String,
            default: "create",
        },
    },
    created: function () {
        this.select();
    },
    watch: {
        tag: function () {
            this.select();
        },
        vault: function () {
            this.select();
        },
    },
    methods: {
        upsert: function (e) {
            // Resets
            e.preventDefault();
            this.error = '';

            // Generates the password if requested
            let self = this;
            if (this.kind > -1) {
                let req = new XMLHttpRequest();
                req.open('POST', '/pwd', false);
                req.setRequestHeader("Content-Type", "application/json");
                req.send(JSON.stringify({
                    length: self.length,
                    digits: self.digits,
                    symbols: self.symbols,
                    no_upper: self.noUpper,
                    repeat: self.repeat,
                }));
                if (req.status === 200) {
                    let res = JSON.parse(req.responseText);
                    self.pass = window.atob(res);
                }
            }

            // Creates or updates the vault
            let req = new XMLHttpRequest()
                , tag = encodeURIComponent(self.tag)
                , vault = encodeURIComponent(self.vault)
            ;
            req.open((self.action === 'create' ? 'POST' : 'PUT'), `/tags/${tag}/vaults/${vault}`, false);
            req.setRequestHeader("Content-Type", "application/json");
            req.send(JSON.stringify({
                user: self.user,
                pass: self.pass,
                url: self.url,
                note: self.note,
            }));
            if (req.status === 200) {
                self.$emit('vault', self.tag, self.vault, 'view');
            } else {
                let res = JSON.parse(req.responseText);
                self.error = res.error;
            }
        },
        select: function() {
            // Resets
            this.resetData();
            if (this.action === 'create') {
                return;
            }
            let self = this
                , req = new XMLHttpRequest()
                , tag = encodeURIComponent(self.tag)
                , vault = encodeURIComponent(self.vault)
            ;

            // Retrieves the vault.
            req.open('GET', `/tags/${tag}/vaults/${vault}`, true);
            req.onload = function() {
                let res = JSON.parse(req.responseText);
                if (this.status === 200) {
                    self.user = res.user;
                    self.url = res.url;
                    self.note = res.note;
                    self.upd_ts = res.upd_ts;
                } else if (this.status === 404) {
                    self.$emit('vault', self.tag, '', 'create');
                } else {
                    self.error = res.error;
                }
            };
            req.send(null);
        },
        resetData: function () {
            Object.assign(this.$data, viewDefaultData());
        }
    },
    template: `<form method="post" class="vault" v-on:submit="upsert">
    <div class="d-flex pt-2 pb-4">
        <div>
            <label class="sr-only" v-bind:for="id + 'Name'">Name</label>
            <input type="text" class="form-control form-control-sm mb-2" v-bind:id="id + 'Name'" name="vault" v-model="vault" placeholder="Vault name" v-if="action === 'create'" minlength="1" required>
            <h2 class="m-0" v-else>\${ vault }</h2>
        </div>
        <div class="ml-auto">
            <button type="submit" class="btn btn-primary btn-sm">Save</button>
        </div>
    </div>
    <div class="alert alert-danger" role="alert" v-if="error.length">\${ error }</div>
    <div class="form-group row">
        <label class="col-sm-3 col-form-label" v-bind:for="id + 'Username'">Username</label>
        <div class="col-sm-9">
            <input type="text" class="form-control" v-bind:id="id + 'Username'" name="user" v-model="user" placeholder="user@example.com" required>
        </div>
    </div>
    <fieldset class="form-group">
        <div class="row">
            <label class="col-sm-3 col-form-label" v-bind:for="id + 'Password'">Password</label>
            <div class="col-sm-9">
                <div class="form-check">
                    <input class="form-check-input" type="radio" name="kind" v-model="kind" v-bind:id="id + 'KindSelfGenerated'" value="0" checked>
                    <label class="form-check-label" v-bind:for="id + 'KindSelfGenerated'">Self-generated</label>
                </div>
                <div class="form-check">
                    <input class="form-check-input" type="radio" name="kind" v-model="kind" v-bind:id="id + 'KindGeneratingCustomized'" value="1">
                    <label class="form-check-label" v-bind:for="id + 'KindGeneratingCustomized'">Generating customized</label>
                </div>
                <div class="form-check">
                    <input class="form-check-input" type="radio" name="kind" v-model="kind" v-bind:id="id + 'KindHandwritten'" value="-1">
                    <label class="form-check-label" v-bind:for="id + 'KindHandwritten'">Handwritten</label>
                </div>
                <hr v-if="kind != 0">
                <div v-show="kind == 1">
                    <div class="form-row">
                        <div class="form-group col-md-4">
                            <label v-bind:for="id + 'Length'">Length</label>
                            <input type="number" class="form-control" name="length" v-model="length" v-bind:id="id + 'Length'" value="64" min="16" max="64">
                        </div>
                        <div class="form-group col-md-4">
                            <label v-bind:for="id + 'Digits'">Number of digits</label>
                            <input type="number" class="form-control" name="digits" v-model="digits" v-bind:id="id + 'Digits'" value="10" min="0" max="64">
                        </div>
                        <div class="form-group col-md-4">
                            <label v-bind:for="id + 'Symbols'">Number of symbols</label>
                            <input type="number" class="form-control" name="symbols" v-model="symbols" v-bind:id="id + 'Symbols'" value="10" min="0" max="64">
                        </div>
                    </div>
                    <div class="form-row">
                        <div class="form-check form-check-inline">
                            <input class="form-check-input" type="checkbox" name="no-upper" v-model="noUpper" v-bind:id="id + 'NoUpper'" value="1">
                            <label class="form-check-label" v-bind:for="id + 'NoUpper'">No uppercase</label>
                        </div>
                        <div class="form-check form-check-inline">
                            <input class="form-check-input" type="checkbox" name="repeat" v-model="repeat" v-bind:id="id + 'Repeat'" value="1" checked>
                            <label class="form-check-label" v-bind:for="id + 'Repeat'">Allow repeat</label>
                        </div>
                    </div>
                </div>
                <div class="form-group" v-show="kind < 0">
                    <input type="password" class="form-control" id="id + 'Password'" name="pass" v-model="pass" placeholder="••••••••" autocomplete="off">
                </div>
            </div>
        </div>
    </fieldset>
    <hr>
    <div class="form-group row">
        <label class="col-sm-3 col-form-label" v-bind:for="id + 'URL'">URL</label>
        <div class="col-sm-9">
            <input type="url" class="form-control" v-bind:id="id + 'URL'" name="url" v-model="url" placeholder="https://example.com" pattern="https://.*">
        </div>
    </div>
    <div class="form-group row">
        <label class="col-sm-3 col-form-label" v-bind:for="id + 'Note'">Note</label>
        <div class="col-sm-9">
            <textarea class="form-control" v-bind:id="id + 'Note'" name="note" v-model="note" rows="3"></textarea>
        </div>
    </div>
    <button type="submit" class="btn btn-primary">Save</button>
</form>`
});