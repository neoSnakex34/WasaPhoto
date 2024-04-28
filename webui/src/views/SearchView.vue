
<script>
    export default{
        data: function(){
            return {
                usersFromQuery: [],
                matchingUsers: [],
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
                this.usersFromQuery = response.data;
                console.log(this.usersFromQuery);

                if (!this.usersFromQuery){
                    alert("No users found");
                } else {
                    this.matchingUsers = this.usersFromQuery.filter(userObj => userObj.user.username.startsWith(this.query.toLowerCase()))
                    console.log(this.matchingUsers);

                    }

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

                    // find user in matching users with banned id
                    let bannedUserFromQuery = this.matchingUsers.find(ufq => ufq.user.userId.identifier === bannedId);
                    bannedUserFromQuery.requestorHasBanned = true;
                    
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

                    // find user in matching users with banned id
                    let ubannedUserFromQuery = this.matchingUsers.find(ufq => ufq.user.userId.identifier === bannedId);
                    ubannedUserFromQuery.requestorHasBanned = false;

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
             v-for="entry in matchingUsers" 
             @mouseenter="entry.showButtons = true"
             @mouseleave="entry.showButtons = false"
             
             style="min-height: 100px;">
            <!-- add href to profile -->
            
            <router-link v-if="!entry.isRequestorBanned && !entry.requestorHasBanned" :to="`/profile/${encodeURIComponent(entry.user.userId.identifier)}`" style="font-size: large;"><strong>{{ entry.user.username }}</strong></router-link>
            <div v-if="entry.isRequestorBanned" class="disabled" style="font-size: large;">{{ entry.user.username }}</div>
            <div v-if="entry.requestorHasBanned" class="disabled" style="font-size: large;">{{ entry.user.username }}</div>

            <div v-if="!entry.isRequestorBanned" class="btn-group"  v-show="entry.showButtons">

                <div v-if="!entry.requestorHasBanned">
                    <button class="btn btn-primary fw-bold rounded-pill ms-auto me-3"  @click="followUser(entry.user.userId.identifier)">Follow</button>
                    <button class="btn btn-danger fw-bold rounded-pill ms-auto me-3" @click="unfollowUser(entry.user.userId.identifier)" >Unfollow</button>
                    <button class="btn btn-secondary fw-bold rounded-pill ms-auto me-3" @click="banUser(entry.user.userId.identifier)">Ban</button>

                </div>
                <button class="btn btn-success fw-bold rounded-pill ms-auto" @click="unbanUser(entry.user.userId.identifier)">Unban</button>

            </div>
        </div>
    </div>
        
</template>