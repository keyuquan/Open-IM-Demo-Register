cd ../

rm open_im_register_api

cd cmd

go build open_im_register_api.go

mv open_im_register_api ../
cd ../

nohup ./open_im_register_api &