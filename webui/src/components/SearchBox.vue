
<script>
    export default{
            props: ['users'],
            data: function(){
                return {
                
                }
            }, 
            methods: {

            },

            computed: {
                removeCurrentUserFromList: function(){
                    return this.users.filter(u => u.username !== this.$root.username) // should i do by id?
                }
            }
        }
    
</script>

<!--- TODO template -->