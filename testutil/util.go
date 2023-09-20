package testutil

import "github.com/cosmos/iavl-bench/bench"

type TreeBuildOptions struct {
	Until       int64
	UntilHash   string
	LoadVersion int64
	Iterator    bench.ChangesetIterator
	Report      func()
	SampleRate  int64
}

func (opts TreeBuildOptions) With10_000() TreeBuildOptions {
	o := &opts
	o.Until = 10_000
	o.UntilHash = "460a9098015aef66f2da7f3d81fedf9a439ea3c3cf61723d535d2d94367858d5"
	return *o
}

func (opts TreeBuildOptions) With25_000() TreeBuildOptions {
	o := &opts
	o.Until = 25_000
	// verified against cosmos/iavl-bench on 2023-09-18
	o.UntilHash = "f1283df353b4766c938d75982c3d69b1eeb7a3c9eea006376ecf7feeab1b9743"
	return *o
}

func (opts TreeBuildOptions) With100_000() TreeBuildOptions {
	o := &opts
	o.Until = 100_000
	o.UntilHash = "e57ab75990453235859416baaccedbaac7b721cd099709ee968321c7822766b1"
	return *o
}

func (opts TreeBuildOptions) With300_000() TreeBuildOptions {
	o := &opts
	o.Until = 300_000
	o.UntilHash = "50a08008a29d76f3502d0a60c9e193a13efa6037a79a9f794652e1f97c2bbc16"
	return *o
}

func (opts TreeBuildOptions) With1_500_000() TreeBuildOptions {
	o := &opts
	o.Until = 1_500_000
	o.UntilHash = "ebc23d2e4e43075bae7ebc1e5db9d5e99acbafaa644b7c710213e109c8592099"
	return *o
}

func NewTreeBuildOptions() TreeBuildOptions {
	var seed int64 = 1234
	var versions int64 = 10_000_000
	bankGen := bench.BankLikeGenerator(seed, versions)
	//bankGen.InitialSize = 10_000
	lockupGen := bench.LockupLikeGenerator(seed, versions)
	//lockupGen.InitialSize = 10_000
	stakingGen := bench.StakingLikeGenerator(seed, versions)
	//stakingGen.InitialSize = 10_000
	itr, err := bench.NewChangesetIterators([]bench.ChangesetGenerator{
		bankGen,
		lockupGen,
		stakingGen,
	})
	if err != nil {
		panic(err)
	}
	opts := TreeBuildOptions{
		Iterator: itr,
	}
	return opts.With25_000()
}

func BankLockup25_000() TreeBuildOptions {
	var seed int64 = 1234
	var versions int64 = 10_000_000
	bankGen := bench.BankLikeGenerator(seed, versions)
	lockupGen := bench.LockupLikeGenerator(seed, versions)
	itr, err := bench.NewChangesetIterators([]bench.ChangesetGenerator{
		bankGen,
		lockupGen,
	})
	if err != nil {
		panic(err)
	}
	opts := TreeBuildOptions{
		Iterator:  itr,
		Until:     25_000,
		UntilHash: "c1dc9dc7d3a8ae025d2a347eea19121e98435b06b421607119bc3cf3cf79be05",
	}
	return opts
}

func BigTreeOptions_100_000() TreeBuildOptions {
	var seed int64 = 1234
	var versions int64 = 200_000
	bankGen := bench.BankLikeGenerator(seed, versions)
	lockupGen := bench.LockupLikeGenerator(seed, versions)
	stakingGen := bench.StakingLikeGenerator(seed, versions)
	itr, err := bench.NewChangesetIterators([]bench.ChangesetGenerator{
		bankGen,
		lockupGen,
		stakingGen,
	})
	if err != nil {
		panic(err)
	}
	opts := TreeBuildOptions{
		Iterator:  itr,
		Until:     10_000,
		UntilHash: "c1dc9dc7d3a8ae025d2a347eea19121e98435b06b421607119bc3cf3cf79be05",
	}
	return opts
}

func BigStartOptions() TreeBuildOptions {
	initialSize := 1_000_000
	var seed int64 = 1234
	var versions int64 = 10_000
	bankGen := bench.BankLikeGenerator(seed, versions)
	bankGen.InitialSize = initialSize
	lockupGen := bench.LockupLikeGenerator(seed, versions)
	lockupGen.InitialSize = initialSize
	stakingGen := bench.StakingLikeGenerator(seed, versions)
	stakingGen.InitialSize = initialSize

	itr, err := bench.NewChangesetIterators([]bench.ChangesetGenerator{
		bankGen,
		lockupGen,
		stakingGen,
	})
	if err != nil {
		panic(err)
	}

	opts := TreeBuildOptions{
		Iterator:  itr,
		Until:     300,
		UntilHash: "b7266b2b30979e1415bcb8ef7fed9637b542213fefd1bb77374aa1f14442aa50", // 300
	}

	return opts
}

func OsmoLike() TreeBuildOptions {
	initialSize := 20_000_000 // revert to 20M!!
	finalSize := int(1.5 * float64(initialSize))
	var seed int64 = 1234
	var versions int64 = 1_000_000
	bankGen := bench.BankLikeGenerator(seed, versions)
	bankGen.InitialSize = initialSize
	bankGen.FinalSize = finalSize
	bankGen2 := bench.BankLikeGenerator(seed+1, versions)
	bankGen2.InitialSize = initialSize
	bankGen2.FinalSize = finalSize
	//lockupGen := bench.LockupLikeGenerator(seed, versions)
	//lockupGen.InitialSize = initialSize
	//stakingGen := bench.StakingLikeGenerator(seed, versions)
	//stakingGen.InitialSize = initialSize

	itr, err := bench.NewChangesetIterators([]bench.ChangesetGenerator{
		bankGen,
		bankGen2,
	})
	if err != nil {
		panic(err)
	}

	opts := TreeBuildOptions{
		Iterator:  itr,
		Until:     10_000,
		UntilHash: "3b43ef49895a7c483ef4b9a84a1f0ddbe7615c9a65bc533f69bc6bf3eb1b3d6c", // 10000
	}

	return opts
}
