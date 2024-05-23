
<script>
    export default{
        data: function(){
            return {
                usersFromQuery: [],
                matchingUsers: [],
                userId: localStorage.getItem('userId'),
                usernames: [],
                query: "",
                // those will be used to show/hide buttons
                // wont be initialized in data for now
                // requestorHasBanned: false,
                // isRequestorBanned: false,
                // requestorHasFollowed: false,
              
            }
        },

        methods: {
            async getUserList(){ // PLEASE NOTE differently from the backend this getUserList func will filter users
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
                    
                    if (this.matchingUsers.length === 0){
                        alert("No users found");
                    }
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

                    // find user in matching users with followed id
                    let followedUserFromQuery = this.matchingUsers.find(ufq => ufq.user.userId.identifier === followedId);
                    followedUserFromQuery.requestorHasFollowed = true;

                } catch(e){
                    if (e.response) {
                        alert(e.response.data)
                    } else {
                        alert(e)
                    }
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
                
                    // find user in matching users with followed id
                    let followedUserFromQuery = this.matchingUsers.find(ufq => ufq.user.userId.identifier === followedId);
                    followedUserFromQuery.requestorHasFollowed = false;
                
                } catch(e){
                   if(e.response){
                        alert(e.response.data)
                    } else {
                        alert(e)
                    }
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
                    // since ban will remove follow
                    bannedUserFromQuery.requestorHasFollowed = false; 
                } catch(e){
                    if(e.response){
                        alert(e.response.data)
                    } else {
                        alert(e)
                    }
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
                    if(e.response){
                        alert(e.response.data)
                    } else {
                        alert(e)
                    }
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
            <button class="btn btn-outline-primary btn-lg fw-bold" type="button" @click="getUserList">Search</button>
        </div>

    </div>
    <div class="container pt-1" style="height: 500px; width: 80%">
        <!-- add a condition in script for excluding banned user -->
        <div class="border-bottom pt-3 pb-3 d-flex justify-content-between align-items-center"
             v-for="entry in matchingUsers" :key="entry.user.userId.identifier" 
             @mouseenter="entry.showButtons = true"
             @mouseleave="entry.showButtons = false"
             
             style="min-height: 100px;">
            <!-- add href to profile -->
            
            <router-link v-if="!entry.isRequestorBanned && !entry.requestorHasBanned" :to="`/profile/${encodeURIComponent(entry.user.userId.identifier)}`" style="font-size: large;"><strong>{{ entry.user.username }}</strong></router-link>
            <div v-if="entry.isRequestorBanned" class="disabled" style="font-size: large;">banned</div>
            <div v-if="entry.requestorHasBanned" class="disabled" style="font-size: large;">{{ entry.user.username }}</div>

            <div v-if="!entry.isRequestorBanned" class="btn-group"  v-show="entry.showButtons">

                <div v-if="!entry.requestorHasBanned">
                    <button v-if="!entry.requestorHasFollowed" class="btn btn-primary fw-bold rounded-pill ms-auto me-3"  @click="followUser(entry.user.userId.identifier)">Follow</button>
                    <button v-if="entry.requestorHasFollowed" class="btn btn-danger fw-bold rounded-pill ms-auto me-3" @click="unfollowUser(entry.user.userId.identifier)" >Unfollow</button>
                    <button class="btn btn-secondary fw-bold rounded-pill ms-auto me-3" @click="banUser(entry.user.userId.identifier)">Ban</button>

                </div>
                <button v-if="entry.requestorHasBanned" class="btn btn-success fw-bold rounded-pill ms-auto" @click="unbanUser(entry.user.userId.identifier)">Unban</button>

            </div>
        </div>
    </div>
        
</template>