Vue.component('tags', {
    delimiters: ['${', '}'],
    data: function () {
        return {
            items: [],
        };
    },
    props: {
        id: {
            type: String,
            required: true
        },
    },
    methods: {
        login: function (e) {
            // Reset
            e.preventDefault();
            this.error = "";
            // Parses
            this.pass = _.trim(this.pass);
            if (this.pass === '') {
                this.error = "missing passphrase";
                return;
            }
            if  (this.signUp && this.pass !== this.confirm) {
                this.error = "unconfirmed passphrase";
                return;
            }
            // Sends
            var vm = this;
            axios.post('/', {
                phrase: vm.pass.toString()
            }).then(function (response) {
                document.location = response.data.goto;
            }).catch(function (error) {
                vm.error = error.response.data.error;
            });
        }
    },
    template: `<div v-bind:id="id">
<div class="clearfix">
    <h6 class="float-left text-white pt-3">Tags</h6>
    <button type="button" class="btn btn-secondary btn-sm my-sm-2 float-right" v-on:click="add">Add</button>
</div>
<ul class="nav flex-column nav-pills flex-md-grow-1" style="overflow: auto">
    <li class="nav-item" v-for="(item, index) in items">
        <p class="nav-link">${ tag }$</p>
    </li>
</ul>
</div>


<form v-bind:action="action" method="post" v-on:submit="login">
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

const tags = new Vue({
    el: '#tags',
});





