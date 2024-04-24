<script>
import Photo from '../components/Photo.vue'

export default {
    components: {
        Photo
    },
    props: ['msg'],
    data: function () {
        return {
        
            profile: {
                username: localStorage.getItem('username'),
                userId: localStorage.getItem('userId'),
                newUsername: '',
                myPhotos: [],
                followerCounter: 0,
                followingCounter: 0,
                photoCounter: 0,
            },
            servedPhotos: [],
            clicked: false,
            deleteToggle: false

        }
    },
    // computed: {
    //     getUserId(){
    //         return this.profile.userId() = localStorage.getItem('userId')
    //     },
    // },
    // this will get the userprofile and sort photos
    async created() {
        await this.getUserProfile()
        await this.updateServedPhotos()

    },



    methods: {

        // generalize it?
        graphicallyLikeBeforeRefresh(index) {
            // on refresh likes will be retrieved from backend
            this.profile.myPhotos[index].likeCounter++
            // will also set liked 
            this.profile.myPhotos[index].likedByCurrentUser = true
        },

        graphicallyUnlikeBeforeRefresh(index) {
            this.profile.myPhotos[index].likeCounter--
            this.profile.myPhotos[index].likedByCurrentUser = false
        },

        toggleDeleteButton() {
            this.deleteToggle = !this.deleteToggle
        },

        // via file selector i call and upload photo
        formUploadSelect() {

            let input = this.$refs.inputForm.files[0]
            if (input) {
                this.uploadPhoto(input)
            }
            
            // empties the selector
            this.$refs.inputForm.value = null

        },


        async updateServedPhotos() {

            if (this.profile.myPhotos === null) {
                this.servedPhotos = []
                return
            }

            let sortedPhotosByDate = this.profile.myPhotos.sort((a, b) => {
                return new Date(b.date) - new Date(a.date)
            })

            let tmpServedPhotos = []
            for (let photo of sortedPhotosByDate) {

                let path = photo.photoPath
                tmpServedPhotos.push(await this.getPhoto(path))
            }

            this.servedPhotos = tmpServedPhotos
        },

        async getUserProfile() {
            try {
                let response = await this.$axios.get(`/users/${this.profile.userId}/profile`
                    , {
                        headers: {
                            Requestor: this.profile.userId
                        }
                    }
                )

                // check for username consistency 
                if (this.profile.username !== response.data.username) {
                    this.profile.username = response.data.username
                    localStorage.setItem('username', response.data.username)
                }

                // populating profile 
                this.profile.followerCounter = response.data.followersCounter
                this.profile.followingCounter = response.data.followingCounter
                this.profile.photoCounter = response.data.photoCounter
                this.profile.myPhotos = response.data.photos
            } catch (e) {
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    alert(e)
                }
            }
        },

        async setMyUsername(newUsername) {
            try {
                let response = await this.$axios.put(`/users/${this.profile.userId}/username`,
                    JSON.stringify(newUsername),
                    {
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    })
                this.profile.username = newUsername // kinda useless but keeps consistency in vue state 
                localStorage.setItem('username', newUsername)
            } catch (e) {

            }
        },

     
        async uploadPhoto(file) {

            let fileReader = new FileReader();

            fileReader.onload = async () => {
                let buffer = fileReader.result
                let photoBinaryU8 = new Uint8Array(buffer)

                try {
                    let response = await this.$axios.post(`/users/${this.profile.userId}/photos`,
                        photoBinaryU8,
                        {
                            headers: {
                                'Content-Type': 'application/octet-stream'
                            }
                        })

                    // alert("Photo uploaded successfully")

                    // this will refresh the profile 
                    await this.getUserProfile()
                    await this.updateServedPhotos()

                } catch (e) {
                    if (e.response.data) {
                        alert(e.response.data)
                    } else {
                        alert(e)
                    }
                }


            }

            fileReader.readAsArrayBuffer(file)

        },

        async deletePhoto(index) {
            let photoId = this.profile.myPhotos[index].photoId.identifier
            
            try {
                let response = await this.$axios.delete(`/users/${this.profile.userId}/photos/${photoId}`,
                    {
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    })

                // this will refresh the profile 
                this.toggleDeleteButton()
                await this.getUserProfile()
                await this.updateServedPhotos()

            } catch (e) {
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                alert(e)
                }
            }
        },
        // THIS WILL CALL SERVEPHOTO IN API 
        async getPhoto(partialPath) {
            let photoId = partialPath.split('/')[1]

            try {
                let response = await this.$axios.get(`users/${this.profile.userId}/photos/${photoId}`,
                    {

                        responseType: 'blob'
                    })

                let servedPhotoUrl = window.URL.createObjectURL(response.data)
                return servedPhotoUrl
            } catch (e) {
                alert(e)
            }
        },
        toggleEditing(newUsername) {
            if (this.clicked) {

                this.setMyUsername(newUsername)
                // this.refresh()
            }
            // alert(this.username)
            this.clicked = !this.clicked
        }
    },

    mounted() {
        this.profile.userId = localStorage.getItem('userId')
        this.getUserProfile()
    }

}
</script>

<template>
    <div class="pt-3 pb-3 border-bottom">

        <div class="d-flex">
            <!-- username and id -->
            <div class="d-flex align-items-baseline ps-1">
                <h3 class="h3"><strong>{{ this.profile.username }}</strong></h3>
                <h6 class="text-muted ms-2">(<strong>{{ this.profile.userId }}</strong>)</h6>
            </div>

            <div class="d-flex ms-auto pe-1">
                <!-- counters -->
                <div class="text-center me-4">
                    <div class="fw-bold">Following</div>
                    <div class="fw-bold">{{ this.profile.followingCounter }}</div>
                </div>
                <div class="text-center ps-4 pe-4 me-4">
                    <div class="fw-bold">Followed</div>
                    <div class="fw-bold">{{ this.profile.followerCounter }}</div>
                </div>
                <div class="text-center">
                    <div class="fw-bold">Photos</div>
                    <div class="fw-bold">{{ this.profile.photoCounter }}</div>
                </div>
            </div>

        </div>



    </div>
    <div class="container pt-3 pb-3 d-flex align-items-center justify-content-center" style="width: 90%;">
        <div v-if="!clicked" class="d-flex input-group align-items-center me-4">
            <input ref="inputForm" class="form-control  rounded-end-0" type="file" accept="image/png, image/jpeg">
            <button class="btn btn-primary  rounded-start-0 fw-bold" type="button"
                @click="formUploadSelect()">Upload</button>

        </div>
        <div class="d-flex">
            <input v-if="clicked" v-model="newUsername" type="text" class="form-control me-4" placeholder="new username"
                style="outline: 2px solid lightcyan;" />
            <button class="btn btn-primary rounded-pill fw-bold" @click="toggleEditing(newUsername)">{{ clicked
                    ? 'Confirm' : 'setUsername' }}</button>
            <button v-if="clicked" class="btn btn-danger rounded-pill  fw-bold ms-2"
                @click="toggleEditing('')">Cancel</button>
        </div>

    </div>


    <div class="border-bottom"></div>

    <!-- photos -->
    <div v-if="!clicked" class="container pt-4 pb-4" style="width: 60%;">
        <div>

            <!--- in stream uploader  will need to be passed from struct -->
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
