Vue.component('sign-form', {
    delimiters: ['${', '}'],
    data: function () {
        return {
            error: '',
            pass: '',
            confirm: ''
        };
    },
    props: {
        action: String,
        id: {
            type: String,
            default: "sign",
        },
        signUp: Boolean,
    },
    methods: {
        login: function (e) {
            // Resets
            e.preventDefault();
            this.reset();

            // Parses
            this.pass = this.pass.trim();
            if (this.pass === '') {
                this.error = "missing passphrase";
                return;
            }
            if  (this.signUp && this.pass !== this.confirm) {
                this.error = "unconfirmed passphrase";
                return;
            }

            // Sends
            let self = this
                , req = new XMLHttpRequest();
            req.open('POST', self.action, false);
            req.setRequestHeader("Content-Type", "application/json");
            req.send(JSON.stringify({phrase: self.pass}));
            let res = JSON.parse(req.responseText);
            if (req.status === 200) {
                document.location = res.goto;
            } else {
                self.error = res.error;
            }
        },
        reset: function () {
            this.error = "";
        }
    },
    template: `<form v-bind:action="action" method="post" v-on:submit="login">
    <div class="alert alert-danger" role="alert" v-if="error.length">\${ error }</div>
    <div class="form-group">
        <label v-bind:for="id + 'Pass'">Passphrase:</label>
        <input v-bind:id="id + 'Pass'" name="pass" v-model="pass" type="password" class="form-control" v-bind:aria-describedby="id + 'PassHelp'" placeholder="One password to rule them all" minlength="16">
        <small v-bind:id="id + 'PassHelp'" class="form-text text-muted">A good passphrase should have at least 16 characters and be difficult to guess (digits, symbols, etc.).</small>
    </div>
    <div class="form-group" v-if="signUp">
        <label v-bind:for="id + 'Confirm'">Confirm passphrase:</label>
        <input v-bind:id="id + 'Confirm'" name="confirm" v-model="confirm" type="password" class="form-control" v-bind:aria-describedby="id + 'ConfirmHelp'" placeholder="Double checking">
        <small v-bind:id="id + 'ConfirmHelp'" class="form-text text-muted">Re-type the passphrase to confirm your choice.</small>
    </div>
    <button type="submit" class="btn btn-primary btn-lg">Submit</button>
</form>`
});


