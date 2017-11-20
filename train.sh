#rm -rf ./downloadedMusic
#mkdir downloadedMusic
source ~/.bashrc
source activate magenta
rm -rf ./pianoroll_rnn_nade
rm /tmp/notesequences.tfrecord
INPUT_DIRECTORY=downloadedMusic

# TFRecord file that will contain NoteSequence protocol buffers.
SEQUENCES_TFRECORD=/tmp/notesequences.tfrecord
convert_dir_to_note_sequences \
  --input_dir=$INPUT_DIRECTORY \
  --output_file=$SEQUENCES_TFRECORD \
  --recursive

pianoroll_rnn_nade_create_dataset \
--input=/tmp/notesequences.tfrecord \
--output_dir=pianoroll_rnn_nade/sequence_examples \
--eval_ratio=0

tensorboard --logdir=pianoroll_rnn_nade/logdir &
pianoroll_rnn_nade_train \
--run_dir=pianoroll_rnn_nade/logdir/run1 \
--sequence_example_file=pianoroll_rnn_nade/sequence_examples/training_pianoroll_tracks.tfrecord \
--hparams="batch_size=24,rnn_layer_sizes=[128, 128, 128]"
