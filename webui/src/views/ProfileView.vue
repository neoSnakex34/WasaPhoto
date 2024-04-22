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
                myPhotosPath: [],
                followerCounter: 0,
                followingCounter: 0,
                photoCounter: 0,
            },

            clicked: false,
            errMsg: null,
        }
    },
    methods: {
        async getUserProfile() {
        
            try {
                let response = await this.$axios.get(`/users/${this.profile.userId}/profile`
                    , {
                        headers: {
                            Authorization: this.profile.userId,
                            Requestor: this.profile.userId
                        }
                    }
                )

                if (this.profile.username !== response.data.username) {
                    this.profile.username = response.data.username
                    localStorage.setItem('username', response.data.username)
                }

                this.profile.followerCounter = response.data.followersCounter
                this.profile.followingCounter = response.data.followingCounter
                this.profile.photoCounter = response.data.photoCounter
                this.profile.myPhotosPath = response.data.photoList
                // alert(this.profile.myPhotosPath)

            } catch (e) {
                this.errMsg = e
                alert(e)
            }
        },

        async setMyUsername(newUsername) {
            try {
                let response = await this.$axios.put(`/users/${this.profile.userId}/username`,
                    JSON.stringify(newUsername),
                    {
                        headers: {
                            Authorization: this.profile.userId,
                            'Content-Type': 'application/json'
                        }
                    })
                this.profile.username = newUsername // kinda useless but keeps consistency in vue state 
                localStorage.setItem('username', newUsername)
            } catch (e) {

            }
        },

        formUploadSelect() {

            let input = this.$refs.inputForm.files[0]
            if (input) {
                this.uploadPhoto(input)
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
                                Authorization: this.profile.userId,
                                'Content-Type': 'application/octet-stream'
                            }
                        })
                } catch (e) {
                    alert(e)
                    // TODO handle error
                }


            }

            fileReader.readAsArrayBuffer(file)

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
    <div class="container pt-4 pb-4" style="width: 60%;">
        <div v-for="url in this.profile.myPhotosPath">
            {{ url }} 
            <Photo
                :src="url"
                uploader="rei"
				date="2021-10-10"
				likes="0"
				liked="false"
            />

        </div>
    </div>
</template>
