<script>
import Photo from '../components/Photo.vue';

export default {
	components: {
		Photo
	},
	data: function() {
		
		return {
			username: localStorage.getItem('username'),
			userId: localStorage.getItem('userId'),
			stream: [], 
			// servedStream: [],
		}
	},
	async mounted() {
		this.getMyStream();
	}, 

	methods: {

			graphicallyLikeBeforeRefresh(id) {

				let photo = this.stream.find(p => p.photoId.identifier === id);
				// on refresh likes will be retrieved from backend
				photo.likeCounter++;
				// will also set liked 
				photo.likedByCurrentUser = true;
			},

			graphicallyUnlikeBeforeRefresh(id) {
				let photo = this.stream.find(p => p.photoId.identifier === id);
				photo.likeCounter--;
				photo.likedByCurrentUser = false;
			},

		   async updateServedStream() {

            if (this.stream === null) {
                this.servedStream = []
                return
            }

            let sortedPhotosByDate = this.stream.sort((a, b) => {
                return new Date(b.date) - new Date(a.date)
            })

            
            for (let photo of sortedPhotosByDate) {
                let path = photo.photoPath
				console.log(photo);
                // array by reference is not a copy
                // cause you know speedy speedy scripting lang
                // doing speedy speedy dumb things
                photo.served = await this.getPhotoForStream(path, photo.uploaderUserId.identifier)

            }

        },

		  // THIS WILL CALL SERVEPHOTO IN API 
        async getPhotoForStream(partialPath, uploaderId) {
            let photoId = partialPath.split('/')[1]

            try {
                let response = await this.$axios.get(`users/${uploaderId}/photos/${photoId}`,
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
	

		async getMyStream() {
			try{
				let response = await this.$axios.get(`/users/${this.userId}/stream`);
				this.stream = response.data;
				console.log(this.stream);
				this.updateServedStream();
			} catch(e) {
				if (e.responde) {
					alert(e.response.data);
				} else {
					alert(e);
				}
			}

		}
	}
	
}
</script>

<template>
	<div>
		
		<!-- <div class="d-flex r input-group pb-4 pt-4 border-bottom"> -->
			<!-- <input type="text" id="findUser" v-model="findUser" class="form-control form-control-lg rounded" placeholder="search"/> -->
  

		<div class="d-flex justify-content-center align-items-center pt-2 pb-2  border-bottom">
			<h1 class="h2"><strong>{{ this.username }}</strong>'s feed</h1>
		</div>

		<div class="container pt-4 pb-4" style="width: 60%;">
			<Photo v-for="photo in stream"

				@like="graphicallyLikeBeforeRefresh(photo.photoId.identifier)"
                @unlike="graphicallyUnlikeBeforeRefresh(photo.photoId.identifier)"

				:src="photo.served" 
                :uploader="photo.uploaderUsername"
                :photoId="photo.photoId.identifier"
                :comments="photo.comments"
                :uploaderId="photo.uploaderUserId.identifier"
                :photoOwnerId="photo.uploaderUserId.identifier"
                :loggedUserId="this.userId"
                :date="photo.date" 
                :likes="photo.likeCounter"
                :liked="photo.likedByCurrentUser" 
				:stream="true"

			 	/>

		</div>

		<!-- <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg> -->
	</div>
</template>

<style>
</style>
