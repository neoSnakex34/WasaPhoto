
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
           
                this.matchingUsers = this.users.filter(user => user.username.startsWith(this.query.toLowerCase())).filter(user => user.userId.identifier !== this.userId);
              
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

        async banUser(bannedId){
              
                try{
                    let response = await this.$axios.put(`/users/${bannedId}/bans/${this.userId}`, {
                        headers: {
                            Authorization: this.userId,
                            'Content-Type': 'application/json'
                        }
                    });
                    // TODO change
                    alert("banned");
                } catch(e){
                    // TODO log errors to alreadyfollowed
                    console.log(e);
                    alert("Error");
                }
        },

        async unbanUser(bannedId){
              
                try{
                    let response = await this.$axios.delete(`/users/${bannedId}/bans/${this.userId}`, {
                        headers: {
                            Authorization: this.userId,
                            'Content-Type': 'application/json'
                        }
                    });
                    // TODO change
                    alert("unbanned");
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
            <button class="btn btn-outline-primary btn-lg fw-bold" type="button" @click="searchUser">Search</button>
        </div>

    </div>
    <div class="container pt-1" style="height: 500px; width: 80%">
        <!-- add a condition in script for excluding banned user -->
        <div class="border-bottom pt-3 pb-3 d-flex justify-content-between align-items-center"
             v-for="user in matchingUsers" 
             @mouseenter="user.showButtons = true"
             @mouseleave="user.showButtons = false"
             
             style="min-height: 100px;">
            <!-- add href to profile -->
            <router-link :to="`/profile/${encodeURIComponent(user.userId.identifier)}`" style="font-size: large;"><strong>{{ user.username }}</strong></router-link>
           
            <div class="btn-group" v-show="user.showButtons">

                <button class="btn btn-primary fw-bold rounded-pill ms-auto me-3"  @click="followUser(user.userId.identifier)">Follow</button>
                <button class="btn btn-danger fw-bold rounded-pill ms-auto me-3" @click="unfollowUser(user.userId.identifier)" >Unfollow</button>
                <button class="btn btn-secondary fw-bold rounded-pill ms-auto me-3" @click="banUser(user.userId.identifier)">Ban</button>
                <button class="btn btn-success fw-bold rounded-pill ms-auto" @click="unbanUser(user.userId.identifier)">Unban</button>

            </div>
        </div>
    </div>
        
</template>