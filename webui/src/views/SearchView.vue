
<script>
    export default{
        data: function(){
            return {
                users: [],
                userId: localStorage.getItem('userId'),
                usernames: [],
                matchingUsers: [],
                query: "",
              
            }
        },

        methods: {
            async searchUser(){
                try{
                    let response = await this.$axios.get("/users", { 
                    headers: {
                        Authorization: this.userId,
                        Requestor: this.userId, 
                    }
                });
                this.users = response.data;
                // this.usernames = this.users.map(user => user.username);
           
                this.matchingUsers = this.users.filter(user => user.username.includes(this.query.toLowerCase())).filter(user => user.userId.identifier !== this.userId);
              
                // this.matchingUsernames = this.usernames.filter(username => username.includes(this.query.toLowerCase()));
                } catch(e){
                    console.log(e);
                    alert("Error");
                }

            }, 
            async followUser(followedId){
                
                try{
                    let response = await this.$axios.put(`/users/${followedId}/followers/${this.userId}`, {
                        headers: {
                            Authorization: this.userId,
                            'Content-Type': 'application/json'
                        }
                    });
                    // TODO change
                    alert("Followed");
                } catch(e){
                    // TODO log errors to alreadyfollowed
                    console.log(e);
                    alert("Error");
                }
            },

            async unfollowUser(followedId){
                try{
                    let response = await this.$axios.delete(`/users/${followedId}/followers/${this.userId}`, {
                        headers: {
                            Authorization: this.userId,
                            'Content-Type': 'application/json'
                        }
                    });
                    // TODO change
                    alert("Unfollowed");
                } catch(e){
                    // TODO log errors to alreadyfollowed
                    console.log(e);
                    alert("Error");
                }

        },
    }
}
    
</script>

<template>
    <div class="d-flex flex-column align-items-center pt-4 pb-4">
        <h3>Search a user</h3>
    </div>
    <div class="container" style="width: 90%;">
        <div class="input-group rounded d-flex align-items-center">
            <input v-model="query" class="form-control form-control-lg" type="text" placeholder="search by username"/>
            <button class="btn btn-success btn-lg fw-bold" type="button" @click="searchUser">Search</button>
        </div>

    </div>
    <div class="container pt-1" style="height: 500px; width: 80%">
        <div class="border-bottom pt-3 pb-3 d-flex justify-content-between align-items-center" v-for="user in matchingUsers" :key="userId">
            <!-- add href to profile -->
            <a>{{ user.username }}</a> <!--  add ref to profile-->
            <div class="btn-group">

                <button class="btn btn-primary fw-bold rounded-pill ms-auto me-3"  @click="followUser(user.userId.identifier)">Follow</button>
                <button class="btn btn-danger fw-bold rounded-pill ms-auto me-3" @click="unfollowUser(user.userId.identifier)" >Unfollow</button>
                <button class="btn btn-warning fw-bold rounded-pill ms-auto">Ban</button>

            </div>
        </div>
    </div>
        
</template>