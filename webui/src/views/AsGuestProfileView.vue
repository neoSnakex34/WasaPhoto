<script>
import Photo from '../components/Photo.vue';

export default {
    props: ['userId'],
    components: {
        Photo
    },

    data: function () {

        return {
            otherProfile: {
                username: "",
                userId: this.userId,
                followingCounter: 0,
                followedCounter: 0,
                photoCounter: 0,
                photos: [],

            },
            backendReadableUserId: decodeURIComponent(this.userId),
            // servedPhotos: [], deprecaated
            bannedByHost: false, // SHOULD BE A PROP
            guest: true // DO NOT CHANGE 
        }
    },

    async created() {
        
        await this.getUserProfileAsGuest()
        await this.updateGuestServedPhotos()
    },

    methods: {

        //  PLEASE NOTE
        // in the sake of simplicity in developing the app for the exam
        // i would not generalize methods with same or similar behaviour into 
        // specific components not to overcomplicate my workflow
        // in a future update this could be generalized
        graphicallyLikeBeforeRefresh(id) {
            let photo = this.otherProfile.photos.find(p => p.photoId.identifier === id);
            photo.likeCounter++;
            photo.likedByCurrentUser = true;
        },

        graphicallyUnlikeBeforeRefresh(id) {
            let photo = this.otherProfile.photos.find(p => p.photoId.identifier === id);
            photo.likeCounter--;
            photo.likedByCurrentUser = false;
        },

        // in a future update this couls be generalized
        // THIS WILL CALL SERVEPHOTO IN API 
        async getPhoto(partialPath) {

            let photoId = partialPath.split('/')[1]

            try {
                let response = await this.$axios.get(`users/${this.otherProfile.userId}/photos/${photoId}`, {
                    responseType: 'blob',
                    headers: {
                        Requestor: localStorage.getItem('userId')
                    }
                })

                let servedPhotoUrl = window.URL.createObjectURL(response.data)
                // alert(servedPhotoUrl)
                return servedPhotoUrl

            } catch (e) {
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    alert(e)
                }
            }
        },

        async updateGuestServedPhotos() {

            if (this.otherProfile.photos === null) {
                this.servedPhotos = []
                return
            }

            let sortedPhotosByDate = this.otherProfile.photos.sort((a, b) => {
                return new Date(b.date) - new Date(a.date)
            })

            for (let photo of sortedPhotosByDate) {
                let path = photo.photoPath
                photo.served = await this.getPhoto(path)
            }

        },

        async getUserProfileAsGuest() {
            try {

                // simulating this with curl should give an error
                let response = await this.$axios.get(`/users/${this.userId}/profile`, {
                    headers: {
                        Requestor: localStorage.getItem('userId')
                    }
                })

                this.otherProfile.username = response.data.username
                // this.otherProfile.userId = response.data.userId.identifier
                this.otherProfile.followedCounter = response.data.followersCounter
                this.otherProfile.followingCounter = response.data.followingCounter
                this.otherProfile.photoCounter = response.data.photoCounter
                this.otherProfile.photos = response.data.photos

            } catch (e) {
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
                <h3 class="h3"><strong>{{ this.otherProfile.username }}</strong></h3>
                <h6 class="text-muted ms-2">(<strong>{{ this.userId }}</strong>)</h6>
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

            <Photo v-for="photo in otherProfile.photos" :key="photo.photoId.identifier"

                @like="graphicallyLikeBeforeRefresh(photo.photoId.identifier)"
                @unlike="graphicallyUnlikeBeforeRefresh(photo.photoId.identifier)"
             
                :src="photo.served" 
                :uploader="this.otherProfile.username"
                :photoId="photo.photoId.identifier"
                :uploaderId="photo.uploaderUserId.identifier"
                :comments="photo.comments"
                :commentsCounter="photo.commentsCounter"
                :photoOwnerId="photo.uploaderUserId.identifier"
                :date="photo.date" 
                :likes="photo.likeCounter"
                :liked="photo.likedByCurrentUser" 
                :guest="this.guest"
                />

        </div>
    </div>
</template>
