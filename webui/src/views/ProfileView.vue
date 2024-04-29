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
                myPhotos: [],
                followerCounter: 0,
                followingCounter: 0,
                photoCounter: 0,
            },
            newUsername: '',
            // servedPhotos: [], deprecated
            clicked: false,
            deleteToggle: false,
            guest: false // DO NOT CHANGE
        }
    },

    async created() {
        if (!this.profile.myPhotos) {
            this.profile.myPhotos = []
        }
        await this.getUserProfile()
        await this.updateServedPhotos()
    },



    methods: {

        //  PLEASE NOTE
        // in the sake of simplicity in developing the app for the exam
        // i would not generalize methods with same or similar behaviour into 
        // specific components not to overcomplicate my workflow
        // in a future update this could be generalized
        graphicallyLikeBeforeRefresh(id) {
            
            // find will always be true and return the right obj cause photoId is unique
            const photo = this.profile.myPhotos.find(p => p.photoId.identifier === id);
            // on refresh likes will be retrieved from backend
            photo.likeCounter++;
            // will also set liked 
            photo.likedByCurrentUser = true;
        },

        graphicallyUnlikeBeforeRefresh(id) {
            const photo = this.profile.myPhotos.find(p => p.photoId.identifier === id);
            photo.likeCounter--;
            photo.likedByCurrentUser = false;
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

            // let tmpServedPhotos = []
            
            // TODO fix as guest
            for (let photo of sortedPhotosByDate) {

                let path = photo.photoPath
                // array by reference is not a copy
                // cause you know speedy speedy scripting lang
                // doing speedy speedy dumb things
                photo.served = await this.getPhoto(path)

            }

            // this.servedPhotos = tmpServedPhotos
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
                if (this.profile.myPhotos === null) {
                    this.profile.myPhotos = []
                }
                console.log(this.profile.myPhotos)
            } catch (e) {
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    if (e.response.status === 403) {
                        alert("You are banned")
                    } else {
                        alert(e)
                    }
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
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    alert(e)
               }
            }

            this.newUsername = '' 
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

                    let newPhoto = response.data
                    this.profile.myPhotos.push(newPhoto)
                    this.profile.photoCounter++ 
                    await this.updateServedPhotos()
                    console.log(this.profile.myPhotos)

                } catch (e) {
                    if (e.response) {
                        alert(e.response.data)
                    } else {
                        alert(e)
                    }
                }


            }

            fileReader.readAsArrayBuffer(file)

        },

        async deletePhoto(id) {
            
            try {
                let response = await this.$axios.delete(`/users/${this.profile.userId}/photos/${id}`,
                    {
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    });

                this.toggleDeleteButton();
              


                // TODO maybe check for errors?
                let idx = this.profile.myPhotos.findIndex(p => p.photoId.identifier === id);
                this.profile.myPhotos.splice(idx, 1)
                this.profile.photoCounter-- 

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
                        responseType: 'blob',
                        headers: {
                            Requestor: localStorage.getItem('userId')
                        }
                    })

                let servedPhotoUrl = window.URL.createObjectURL(response.data)
                return servedPhotoUrl
            } catch (e) {
                alert(e)
            }
        },

        toggleEditing(newUsername) {
            if ((this.clicked) && (newUsername !== "")){

                this.setMyUsername(newUsername)
            }
            // alert(this.username)
            this.clicked = !this.clicked
        }
    },



}
</script>

<template>
    <!--- could be generalized into a component need to decide if i want to... -->
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
                    <div class="fw-bold">Followed by</div>
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
            <input v-if="clicked" v-model="newUsername" ref="username-form" type="text" class="form-control me-4" placeholder="new username"
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
            <Photo v-for="photo in profile.myPhotos" :key="photo.photoId.identifier"
                
                @like="graphicallyLikeBeforeRefresh(photo.photoId.identifier)"
                @unlike="graphicallyUnlikeBeforeRefresh(photo.photoId.identifier)" 
                @toggle-delete="this.deleteToggle = !this.deleteToggle"
                @delete-event="deletePhoto(photo.photoId.identifier)"
         
                
                :src="photo.served" 
                :uploader="this.profile.username"
                :photoId="photo.photoId.identifier"
                :comments="photo.comments"
                :uploaderId="photo.uploaderUserId.identifier"
                :photoOwnerId="photo.uploaderUserId.identifier"
                :loggedUserId="this.profile.userId"
                :date="photo.date" 
                :likes="photo.likeCounter"
                :liked="photo.likedByCurrentUser" 
                :delete = this.deleteToggle
                :guest = this.guest

                />

        </div>
    </div>
</template>
