<script>
import Photo from '../components/Photo.vue';

export default {
    props : ['userId'],
    components: {
        Photo
    },

    data: function () {
       
        return {
            otherProfile: {
                username: "",
                
                followingCounter: 0,
                followedCounter: 0,
                photoCounter: 0,
                photos: [],

            },
            backendReadableUserId: decodeURIComponent(this.userId),
            servedPhotos: [],
            bannedByHost: false, // SHOULD BE A PROP

        }
    }, 

    async mounted(){
       
        await this.getUserProfileAsGues()
    }, 

    methods:{

        async getUserProfileAsGues() {
            try {

                // TODO check if ban works 
                let response = await this.$axios.get(`/users/${this.userId}/profile`, {
                    headers: {
                        Requestor: localStorage.getItem('userId')
                    }
                })

                
                this.otherProfile.username = response.data.username
                this.otherProfile.followedCounter = response.data.followersCounter
                this.otherProfile.followingCounter = response.data.followingCounter


            }catch(e){
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    alert(e)
                }
            }
        }

    }, 

}
</script>

<template>
    <div class="pt-3 pb-3 border-bottom">

        <div class="d-flex">
            <!-- username and id -->
            <div class="d-flex align-items-baseline ps-1">
                <h3 class="h3"><strong>{{this.otherProfile.username}}</strong></h3>
                <h6 class="text-muted ms-2">(<strong>{{this.userId}}</strong>)</h6>
            </div>

            <div class="d-flex ms-auto pe-1">
                <!-- counters -->
                <div class="text-center me-4">
                    <div class="fw-bold">Following</div>
                    <div class="fw-bold">{{ this.otherProfile.followingCounter }}</div>
                </div>
                <div class="text-center ps-4 pe-4 me-4">
                    <div class="fw-bold">Followed by</div>
                    <div class="fw-bold">{{ this.otherProfile.followedCounter }}</div>
                </div>
                <div class="text-center">
                    <div class="fw-bold">Photos</div>
                    <div class="fw-bold">{{ this.otherProfile.photoCounter }}</div>
                </div>
            </div>

        </div>



    </div>

    <div class="border-bottom"></div>

    <!-- photos v if not banned -->
    <div class="container pt-4 pb-4" style="width: 60%;">
        <div v-if="!bannedByHost">   

            <Photo v-for="(photo, index) in servedPhotos"
                
                @like="graphicallyLikeBeforeRefresh(index)"
                @unlike="graphicallyUnlikeBeforeRefresh(index)" 
                @toggle-delete="this.deleteToggle = !this.deleteToggle"
                @delete-event="deletePhoto(index)"
                :src="photo" :uploader="this.profile.username"
                :photoId="this.profile.myPhotos[index].photoId.identifier"
                :uploaderId="this.profile.myPhotos[index].uploaderUserId.identifier"
                :date="this.profile.myPhotos[index].date" :likes="this.profile.myPhotos[index].likeCounter"
                :liked="this.profile.myPhotos[index].likedByCurrentUser" 
                :delete = this.deleteToggle
                />

        </div>
    </div>
</template>
