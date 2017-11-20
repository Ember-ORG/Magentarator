#rm -rf ./downloadedMusic
#mkdir downloadedMusic
pkill -9 pianoroll_rnn_nade
pkill -9 gnome-terminal
source ~/.bashrc
source activate magenta
pianoroll_rnn_nade_generate \
--run_dir=./pianoroll_rnn_nade/logdir/run1 \
--output_dir=./pianoroll_rnn_nade/generated \
--num_outputs=3 \
--num_steps=2048 \
--primer_pitches="[67]" \
--hparams="batch_size=24,rnn_layer_sizes=[128, 128, 128]"
