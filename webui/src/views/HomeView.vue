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
			stream: [], // TODO
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {

			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();

			}
			this.loading = false;
		},
		
		doLogout() {
			localStorage.removeItem('userId');
			localStorage.removeItem('username');
			this.$router.push('/login');
		},
		

	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		
		<!-- <div class="d-flex r input-group pb-4 pt-4 border-bottom"> -->
			<!-- <input type="text" id="findUser" v-model="findUser" class="form-control form-control-lg rounded" placeholder="search"/> -->
  			<div class="container pb-4 pt-4 border-bottom" style="width: 70%">
				<div class="d-flex input-group align-items-center">
  					<input class="form-control form-control-lg rounded" id="formFile" type="file" accept="image/png, image/jpeg">
				
					<button class="btn btn-primary btn-lg rounded-pill fw-bold ms-2" type="button" id="button-addon2">Upload</button>
				</div>
			</div>
		

		<div class="d-flex justify-content-center align-items-center pt-2 pb-2  border-bottom">
			<h1 class="h2"><strong>{{ this.username }}</strong>'s feed</h1>
		</div>

		<div class="container pt-4 pb-4" style="width: 60%;">
			<Photo
				:src="'rei.jpg'"
				uploader="rei"
				date="2021-10-10"
				likes="0"
				liked="false"
			 	/>
		
		</div>

		<!-- <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg> -->
	</div>
</template>

<style>
</style>
